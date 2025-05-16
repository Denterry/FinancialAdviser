from __future__ import annotations

import asyncio
import logging
import signal
import sys
from typing import NoReturn

from core.db import init_pool, close_pool
from grpc.server import serve_grpc

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s | %(levelname)-8s | %(name)s | %(message)s",
    stream=sys.stdout,
)
logger = logging.getLogger("brain.main")


async def _runner() -> NoReturn:
    # 1. DB-pool
    await init_pool()
    logger.info("PostgreSQL pool initialised")

    # 2. gRPC-сервер как отдельная задача
    grpc_task = asyncio.create_task(serve_grpc(), name="grpc-server")

    # 3. graceful shutdown через Event + signal-handlers
    stop_event = asyncio.Event()

    def _graceful_stop(sig: signal.Signals) -> None:  # noqa: D401
        logger.warning("Received %s → shutting down …", sig.name)
        stop_event.set()

    loop = asyncio.get_running_loop()
    for s in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(s, _graceful_stop, s)

    # 4. ожидаем сигнала
    await stop_event.wait()

    # 5. завершаем gRPC-сервер
    grpc_task.cancel()
    try:
        await grpc_task
    except asyncio.CancelledError:
        logger.info("gRPC task cancelled")

    # 6. закрываем пул БД
    await close_pool()
    logger.info("DB pool closed")

    logger.info("Brain-service stopped gracefully")
    # asyncio.run не возвращает управление; просто «выпадем»


if __name__ == "__main__":
    try:
        asyncio.run(_runner())
    except KeyboardInterrupt:
        pass
