from __future__ import annotations

import asyncpg
import os
from typing import Optional


PG_HOST = os.getenv("POSTGRES_HOST", "postgres")
PG_PORT = int(os.getenv("POSTGRES_PORT", "5432"))
PG_USER = os.getenv("POSTGRES_USER", "postgres")
PG_PASSWORD = os.getenv("POSTGRES_PASSWORD", "postgres")
PG_DB = os.getenv("POSTGRES_DB", "brain")
PG_SSL = os.getenv("POSTGRES_SSL_MODE", "disable")  # disable / require / verify-full

POOL_MIN_SIZE = int(os.getenv("PG_POOL_MIN", "2"))
POOL_MAX_SIZE = int(os.getenv("PG_POOL_MAX", "10"))

_pool: Optional[asyncpg.Pool] = None


async def init_pool() -> asyncpg.Pool:
    """
    Создаёт пул, если он ещё не создан.  
    Вызывать один раз при запуске сервиса
    """
    global _pool
    if _pool is None:
        dsn = (
            f"postgresql://{PG_USER}:{PG_PASSWORD}@{PG_HOST}:{PG_PORT}/{PG_DB}"
            f"?sslmode={PG_SSL}"
        )
        _pool = await asyncpg.create_pool(
            dsn=dsn,
            min_size=POOL_MIN_SIZE,
            max_size=POOL_MAX_SIZE,
            command_timeout=60,
        )
    return _pool


async def get_pool() -> asyncpg.Pool:
    """
    Безопасно возвращает пул.  
    Если init_pool ещё не вызывался — кидает RuntimeError
    """
    if _pool is None:
        raise RuntimeError("Database pool is not initialized; call init_pool() first.")
    return _pool


async def close_pool() -> None:
    """
    Корректно закрывает пул
    """
    global _pool
    if _pool is not None:
        await _pool.close()
        _pool = None
