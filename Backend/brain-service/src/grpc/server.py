from __future__ import annotations

import asyncio
import logging
import os
import uuid
from typing import AsyncIterator, Optional

import grpc
from asyncpg.pool import Pool

from proto import brain_pb2 as pb
from proto import brain_pb2_grpc as pb_grpc

from core.db import init_pool, close_pool
from repo.chats import ChatRepository
from repo.messages import MessageRepository
from services.llm_orchestrator import LLMOrchestrator
from services.ml_client import MLClient


def _chat_dict_to_pb(row: dict) -> pb.Chat:
    return pb.Chat(
        id=str(row["id"]),
        user_id=str(row["user_id"]),
        title=row["title"],
        created_at=row["created_at"].isoformat(),
        updated_at=row["updated_at"].isoformat(),
    )


def _msg_dict_to_pb(row: dict) -> pb.Message:
    role_map = {
        "user": pb.Role.ROLE_USER,
        "assistant": pb.Role.ROLE_ASSISTANT,
        "system": pb.Role.ROLE_SYSTEM,
    }
    return pb.Message(
        id=row["id"],
        chat_id=str(row["chat_id"]),
        role=role_map.get(row["role"], pb.Role.ROLE_UNSPECIFIED),
        content=row["content"],
        token_count=row["token_count"],
        created_at=row["created_at"].isoformat(),
    )


class BrainServiceServicer(pb_grpc.BrainServiceServicer):
    """Имплементация RPC-методов из brain.proto."""

    def __init__(
        self,
        pool: Pool,
        chat_repo: ChatRepository,
        msg_repo: MessageRepository,
        orchestrator: LLMOrchestrator,
    ):
        self.pool = pool
        self.chats = chat_repo
        self.msgs = msg_repo
        self.orch = orchestrator
        self.log = logging.getLogger("brain.grpc")

    async def CreateChat(
        self, request: pb.CreateChatRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb.CreateChatResponse:
        row = await self.chats.create_chat(uuid.UUID(request.user_id),
                                           request.title)
        return pb.CreateChatResponse(chat=_chat_dict_to_pb(row))

    async def ListChats(
        self, request: pb.ListChatsRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb.ListChatsResponse:
        limit = request.page_size or 20
        cursor: Optional[uuid.UUID] = (
            uuid.UUID(request.page_token) if request.page_token else None
        )
        rows, next_cursor = await self.chats.list_chats(
            uuid.UUID(request.user_id), limit, cursor
        )
        return pb.ListChatsResponse(
            chats=[_chat_dict_to_pb(r) for r in rows],
            next_page_token=str(next_cursor) if next_cursor else "",
        )

    async def DeleteChat(
        self, request: pb.DeleteChatRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb.DeleteChatResponse:
        await self.chats.delete_chat(
            uuid.UUID(request.user_id), uuid.UUID(request.chat_id)
        )
        return pb.DeleteChatResponse()

    async def ListMessages(
        self, request: pb.ListMessagesRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb.ListMessagesResponse:
        limit = request.page_size or 50
        cursor: Optional[int] = int(request.page_token) \
            if request.page_token else None
        rows, next_cursor = await self.msgs.list_messages(
            uuid.UUID(request.chat_id), limit, cursor
        )
        return pb.ListMessagesResponse(
            messages=[_msg_dict_to_pb(r) for r in rows],
            next_page_token=str(next_cursor) if next_cursor else "",
        )

    async def StreamMessage(
        self, request: pb.StreamMessageRequest,
        context: grpc.aio.ServicerContext,
    ) -> AsyncIterator[pb.StreamMessageResponse]:
        """
        • Сохраняем пользовательское сообщение.
        • Получаем асинхронный генератор чанков от LLMOrchestrator.
        • Стримим их в клиента.
        """
        # 1) cоздать чат «на лету», если отсутствует
        chat_id = (
            uuid.UUID(request.chat_id)
            if request.chat_id
            else uuid.UUID(
                (
                    await self.chats.create_chat(
                        uuid.UUID(request.user_id),
                        "New chat",
                    )
                )["id"]
            )
        )

        # 2) cохранить входящее сообщение
        await self.msgs.create_message(
            chat_id=chat_id,
            role="user",
            content=request.content,
            token_count=0,  # посчитаем позже, если нужно
        )

        # 3) получить потоковых ассистент-чанков
        async for chunk, is_final, tokens in self.orch.stream_response(
            chat_id, request.content
        ):
            yield pb.StreamMessageResponse(
                content_chunk=chunk, is_final=is_final, tokens_used=tokens
            )

    async def Ping(
        self, request: pb.PingRequest, context: grpc.aio.ServicerContext
    ) -> pb.PingResponse:
        return pb.PingResponse(msg="pong")


_GRPC_PORT = int(os.getenv("GRPC_PORT", "50052"))


async def serve_grpc() -> None:
    logging.basicConfig(
        level=logging.INFO, format="%(asctime)s  %(levelname)s  %(msg)s",
    )

    # 1. DB-пул
    pool = await init_pool()

    # 2. репозитории & сервисы
    chat_repo = ChatRepository(pool)
    msg_repo = MessageRepository(pool)
    ml_client = MLClient()
    orchestrator = LLMOrchestrator(msg_repo, ml_client)

    # 3. gRPC-сервер
    server = grpc.aio.server(options=[("grpc.max_send_message_length",
                                       16 * 1024 * 1024)])
    pb_grpc.add_BrainServiceServicer_to_server(
        BrainServiceServicer(pool, chat_repo, msg_repo, orchestrator), server
    )
    listen_addr = f"[::]:{_GRPC_PORT}"
    server.add_insecure_port(listen_addr)
    await server.start()
    logging.info("Brain gRPC-server started on %s", listen_addr)
    await server.wait_for_termination()


if __name__ == "__main__":
    try:
        asyncio.run(serve_grpc())
    finally:
        asyncio.run(close_pool())
