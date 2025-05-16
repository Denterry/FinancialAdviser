from __future__ import annotations

import asyncpg
from typing import List, Optional, Tuple, Dict
from uuid import UUID


class MessageRepository:
    def __init__(self, pool: asyncpg.Pool):
        self.pool = pool

    async def list_messages(
        self,
        chat_id: UUID,
        limit: int = 50,
        cursor: Optional[int] = None,
    ) -> Tuple[List[Dict], Optional[int]]:
        """
        Возвращает сообщения чата в хрон. порядке (старые → новые).
        """
        if cursor:
            rows = await self.pool.fetch(
                """
                SELECT id, chat_id, role, content,
                       token_count, created_at
                  FROM messages
                 WHERE chat_id = $1
                   AND id < $2
              ORDER BY id DESC
                 LIMIT $3
                """,
                chat_id, cursor, limit,
            )
        else:
            rows = await self.pool.fetch(
                """
                SELECT id, chat_id, role, content,
                       token_count, created_at
                  FROM messages
                 WHERE chat_id = $1
              ORDER BY id DESC
                 LIMIT $2
                """,
                chat_id, limit,
            )

        # Переворачиваем, чтобы отдать «снизу вверх» (старые → новые)
        rows = list(rows)[::-1]
        next_cursor = rows[0]["id"] if rows else None
        return [dict(r) for r in rows], next_cursor

    async def create_message(
        self,
        chat_id: UUID,
        role: str,
        content: str,
        token_count: int,
    ) -> Dict:
        """
        Сохраняет сообщение (user/assistant/system) и возвращает dict.
        """
        row = await self.pool.fetchrow(
            """
            INSERT INTO messages (chat_id, role, content, token_count)
            VALUES ($1, $2, $3, $4)
            RETURNING id, chat_id, role, content,
                      token_count, created_at
            """,
            chat_id, role, content, token_count,
        )
        return dict(row)
