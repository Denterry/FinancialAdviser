from __future__ import annotations

import logging
import os
import time
from typing import AsyncIterator, List, Tuple
from uuid import UUID

import openai
from openai import AsyncOpenAI

from repo.messages import MessageRepository
from services.ml_client import MLClient
from utils.tokenizer import rough_token_count

from utils.prompt_builder import (
    build_chat_prompt,
    append_disclaimer,
)

logger = logging.getLogger("brain.orchestrator")

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")
OPENAI_MODEL = os.getenv("OPENAI_MODEL", "gpt-3.5-turbo")
OPENAI_TEMP = float(os.getenv("OPENAI_TEMPERATURE", "0.7"))
HISTORY_LIMIT = int(os.getenv("CHAT_HISTORY_LIMIT", "20"))


class LLMOrchestrator:
    def __init__(
        self,
        msg_repo: MessageRepository,
        ml_client: MLClient,
    ):
        self.msg_repo = msg_repo
        self.ml = ml_client
        self.openai = AsyncOpenAI(api_key=OPENAI_API_KEY)

    async def stream_response(
        self, chat_id: UUID, user_message: str
    ) -> AsyncIterator[Tuple[str, bool, int]]:
        """
        Генерирует (chunk, is_final, total_tokens)
        """
        # 1. история чата
        history_rows, _ = await self.msg_repo.list_messages(
            chat_id, limit=HISTORY_LIMIT
        )
        history_text = [
            r["content"] for r in history_rows if r["role"] != "system"
        ]

        # 2. динамический sentiment
        try:
            senti = await self.ml.analyze_sentiment(user_message)
            sentiment_line = f"Detected market sentiment: {senti['summary']}."
        except Exception as e:
            logger.warning("ML-service sentiment failed: %s", e)
            sentiment_line = "Market sentiment: unavailable."

        context_lines: List[str] = [sentiment_line]

        # 3. финальный prompt
        openai_messages = build_chat_prompt(
            user_question=user_message,
            context_lines=context_lines + history_text,  # история ⊕ контекст
        )

        start_ts = time.perf_counter()
        assistant_chunks: List[str] = []
        total_tokens = 0

        try:
            stream = await self.openai.chat.completions.create(
                model=OPENAI_MODEL,
                messages=openai_messages,
                temperature=OPENAI_TEMP,
                stream=True,
            )
            async for part in stream:
                delta = part.choices[0].delta.content or ""
                if delta:
                    assistant_chunks.append(delta)
                    total_tokens += rough_token_count(delta)
                    yield delta, False, 0
        except openai.OpenAIError as e:
            err_msg = f"[LLM error] {str(e)}"
            logger.error(err_msg)
            yield err_msg, True, 0
            return

        latency_ms = int((time.perf_counter() - start_ts) * 1000)
        full_answer = append_disclaimer("".join(assistant_chunks))

        await self.msg_repo.create_message(
            chat_id=chat_id,
            role="assistant",
            content=full_answer,
            token_count=total_tokens,
        )

        yield "", True, total_tokens
        logger.info("chat=%s tokens=%d latency=%dms",
                    chat_id, total_tokens, latency_ms)
