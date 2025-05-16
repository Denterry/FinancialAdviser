from fastapi import APIRouter

from src.api.v1.endpoints import sentiment, trends, trading

api_router = APIRouter()

api_router.include_router(
    sentiment.router,
    prefix="/sentiment",
    tags=["sentiment"],
)

api_router.include_router(
    trends.router,
    prefix="/trends",
    tags=["trends"],
)

api_router.include_router(
    trading.router,
    prefix="/trading",
    tags=["trading"],
) 