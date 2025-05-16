from __future__ import annotations

import asyncpg
from typing import List, Optional, Tuple, Dict
from uuid import UUID


class ChatRepository:
    def __init__(self, pool: asyncpg.Pool):
        self.pool = pool

    async def create_chat(self, user_id: UUID, title: str) -> Dict:
        """
        Создаёт новый чат и возвращает строку в виде dict.
        """
        row = await self.pool.fetchrow(
            """
            INSERT INTO chats (user_id, title)
            VALUES ($1, $2)
            RETURNING id, user_id, title,
                     created_at, updated_at
            """,
            user_id, title,
        )
        return dict(row)

    async def list_chats(
        self,
        user_id: UUID,
        limit: int = 20,
        cursor: Optional[UUID] = None,
    ) -> Tuple[List[Dict], Optional[UUID]]:
        """
        Возвращает список чатов пользователя (последние сверху)
        и курсор на следующую страницу.
        """
        if cursor:
            rows = await self.pool.fetch(
                """
                SELECT id, user_id, title, created_at, updated_at
                  FROM chats
                 WHERE user_id = $1
                   AND id < $2
              ORDER BY created_at DESC
                 LIMIT $3
                """,
                user_id, cursor, limit,
            )
        else:
            rows = await self.pool.fetch(
                """
                SELECT id, user_id, title, created_at, updated_at
                  FROM chats
                 WHERE user_id = $1
              ORDER BY created_at DESC
                 LIMIT $2
                """,
                user_id, limit,
            )

        next_cursor = rows[-1]["id"] if rows else None
        return [dict(r) for r in rows], next_cursor

    async def delete_chat(self, user_id: UUID, chat_id: UUID) -> None:
        """
        Мягкое удаление не нужно — просто удаляем чат
        (ON DELETE CASCADE уберёт сообщения).
        """
        await self.pool.execute(
            """
            DELETE FROM chats
             WHERE id = $1
               AND user_id = $2
            """,
            chat_id, user_id,
        )
