from __future__ import annotations

import logging
import os
from typing import Dict, List

from grpc import aio

from proto.ml.v1 import ml_pb2 as ml_pb
from proto.ml.v1 import ml_pb2_grpc as ml_grpc

logger = logging.getLogger("brain.ml_client")

ML_HOST = os.getenv("ML_SERVICE_HOST", "ml-service")
ML_PORT = int(os.getenv("ML_SERVICE_PORT", "50053"))
_ML_ENDPOINT = f"{ML_HOST}:{ML_PORT}"


class MLClient:
    """Lazy-инициализируемый aio-клиент к ml-service."""

    def __init__(self) -> None:
        self._channel: aio.Channel | None = None
        self._stub: ml_grpc.MLServiceStub | None = None

    async def _ensure_channel(self) -> ml_grpc.MLServiceStub:
        if self._stub is None:
            logger.info("Connecting to ML-service at %s …", _ML_ENDPOINT)
            self._channel = aio.insecure_channel(
                _ML_ENDPOINT,
                options=[
                    ("grpc.keepalive_time_ms", 30_000),
                    ("grpc.keepalive_timeout_ms", 10_000),
                ],
            )
            self._stub = ml_grpc.MLServiceStub(self._channel)
        return self._stub

    async def analyze_sentiment(self, text: str) -> Dict:
        stub = await self._ensure_channel()
        req = ml_pb.SentimentAnalysisRequest(text=text)
        resp = await stub.AnalyzeSentiment(req, timeout=5)  # сек
        # агрегируем в простой dict
        if not resp.results:
            return {"summary": "neutral / no data", "details": []}

        # минимальный summary: самый высокий confidence
        top = max(resp.results, key=lambda r: r.confidence)
        summary = f"{top.ticker or 'General'} → {top.sentiment} \
            ({top.confidence:.2f})"
        details = [
            {
                "ticker": r.ticker,
                "sentiment": r.sentiment,
                "confidence": r.confidence,
            }
            for r in resp.results
        ]
        return {"summary": summary, "details": details}

    async def forecast_price(self, ticker: str, horizon_days: int = 7) -> Dict:
        stub = await self._ensure_channel()
        req = ml_pb.PriceForecastRequest(ticker=ticker,
                                         horizon_days=horizon_days)
        resp = await stub.ForecastPrice(req, timeout=5)
        f = resp.forecast
        return {
            "ticker": f.ticker,
            "predicted_price": f.predicted_price,
            "change_pct": f.expected_change_pct,
            "confidence": f.confidence,
            "model": f.model,
        }

    async def forecast_prices(
        self, tickers: List[str], horizon_days: int = 7
    ) -> List[Dict]:
        stub = await self._ensure_channel()
        req = ml_pb.PriceForecastBatchRequest(
            tickers=tickers, horizon_days=horizon_days
        )
        resp = await stub.ForecastPrices(req, timeout=8)
        return [
            {
                "ticker": f.ticker,
                "predicted_price": f.predicted_price,
                "change_pct": f.expected_change_pct,
                "confidence": f.confidence,
                "model": f.model,
            }
            for f in resp.forecasts
        ]

    async def close(self) -> None:
        if self._channel:
            await self._channel.close()
            logger.info("ML-service channel closed")
