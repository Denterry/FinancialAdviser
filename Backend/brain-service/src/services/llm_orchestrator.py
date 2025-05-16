import asyncio
from typing import Optional, Tuple, List

from utils.ticker_extractor import extract_tickers_from_text
from utils.prompt_builder import build_chat_prompt, append_disclaimer
from utils.llm_client import call_openai
from services.data_fetcher import (
    fetch_sentiment_summary,
    fetch_price_prediction,
)
from db.chat_storage import ChatHistory


class ChatOrchestrator:
    """
    Orchestrates the full lifecycle of processing a chat message:
    - Identifies tickers
    - Gathers context from other services
    - Builds a prompt
    - Calls the LLM
    - Stores chat history
    """

    def __init__(self):
        self.chat_store = ChatHistory()

    async def process_user_query(
        self,
        user_id: str,
        message: str,
        chat_id: Optional[int] = None
    ) -> Tuple[str, int]:
        # Start or resume a chat
        chat_id = await self.chat_store.start_or_continue_chat(
            user_id,
            chat_id,
        )

        tickers = extract_tickers_from_text(message)
        context_lines = await self._gather_context_lines(tickers)

        prompt_messages = build_chat_prompt(message, context_lines)

        llm_raw_response = await call_openai(prompt_messages)
        llm_final_response = append_disclaimer(llm_raw_response)

        await self.chat_store.store_message(
            chat_id,
            role="user",
            content=message,
        )
        await self.chat_store.store_message(
            chat_id,
            role="assistant",
            content=llm_final_response,
        )

        return llm_final_response, chat_id

    async def _gather_context_lines(self, tickers: List[str]) -> List[str]:
        """
        Concurrently fetch sentiment and prediction for each ticker
        and return a list of enriched prompt lines.
        """
        if not tickers:
            return []

        tasks = []
        for ticker in tickers:
            tasks.append(fetch_sentiment_summary(ticker))
            tasks.append(fetch_price_prediction(ticker))

        results = await asyncio.gather(*tasks, return_exceptions=True)

        context = []
        for i in range(0, len(results), 2):
            sentiment_result = results[i]
            prediction_result = results[i + 1]
            ticker = tickers[i // 2]

            if not isinstance(sentiment_result, Exception):
                context.append(f"Sentiment on {ticker}: {sentiment_result}")

            if not isinstance(prediction_result, Exception):
                context.append(f"Prediction for {ticker}: {prediction_result}")

        return context
