import httpx
import logging
from typing import Dict, List

logger = logging.getLogger(__name__)

X_SERVICE_URL = "http://x-service:8080"
ML_SERVICE_URL = "http://ml-service:8080"


async def fetch_sentiment_summary(ticker: str) -> Dict:
    try:
        async with httpx.AsyncClient(timeout=5.0) as client:
            res = await client.get(
                f"{X_SERVICE_URL}/api/v1/sentiment/{ticker}/summary",
            )
            res.raise_for_status()
            return res.json()
    except Exception as e:
        logger.warning(
            f"[context] Failed to fetch sentiment for {ticker}: {e}",
        )
        return {"sentiment": "Unknown"}


async def fetch_latest_price(ticker: str) -> Dict:
    try:
        async with httpx.AsyncClient(timeout=5.0) as client:
            res = await client.get(
                f"{X_SERVICE_URL}/api/v1/trading/{ticker}/latest",
            )
            res.raise_for_status()
            return res.json()
    except Exception as e:
        logger.warning(f"[context] Failed to fetch price for {ticker}: {e}")
        return {"current_price": "N/A", "price_change_24h": "N/A"}


async def fetch_price_prediction(ticker: str) -> Dict:
    try:
        async with httpx.AsyncClient(timeout=5.0) as client:
            res = await client.get(f"{ML_SERVICE_URL}/api/v1/predict/{ticker}")
            res.raise_for_status()
            return res.json()
    except Exception as e:
        logger.warning(
            f"[context] Failed to fetch prediction for {ticker}: {e}",
        )
        return {"prediction": "N/A"}


async def gather_context_for_tickers(tickers: List[str]) -> Dict[str, Dict]:
    """
    Gathers sentiment, price, and prediction for each ticker.
    Returns:
        {
            "TSLA": {
                "sentiment": "65% bullish",
                "current_price": 723.54,
                "price_change_24h": 1.23,
                "prediction": "moderate uptrend"
            },
            ...
        }
    """
    import asyncio

    context = {}

    async def fetch_all(ticker: str):
        sentiment = await fetch_sentiment_summary(ticker)
        price = await fetch_latest_price(ticker)
        prediction = await fetch_price_prediction(ticker)

        context[ticker] = {
            "sentiment": sentiment.get("summary", "Unknown"),
            "current_price": price.get("close", "N/A"),
            "price_change_24h": price.get("change_percent", "N/A"),
            "prediction": prediction.get("trend", "N/A"),
        }

    await asyncio.gather(*(fetch_all(t) for t in tickers))
    return context
