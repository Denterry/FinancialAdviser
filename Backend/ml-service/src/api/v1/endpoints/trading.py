from typing import Dict, Any, List
from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel

from src.ml.models import TradingAgent
from src.core.config import settings

router = APIRouter()


class MarketData(BaseModel):
    symbol: str
    price: float
    volume: float
    timestamp: str
    indicators: Dict[str, float]


class TradingRequest(BaseModel):
    market_data: MarketData
    sentiment_data: Dict[str, float]
    trend_data: Dict[str, List[float]]


class TradingResponse(BaseModel):
    action: str
    confidence: float
    parameters: Dict[str, Any]
    explanation: str


def get_trading_agent() -> TradingAgent:
    if not settings.TRADING_ENABLED:
        raise HTTPException(
            status_code=403,
            detail="Trading is not enabled",
        )
    return TradingAgent()


@router.post("/analyze", response_model=TradingResponse)
async def analyze_trading_opportunity(
    request: TradingRequest,
    agent: TradingAgent = Depends(get_trading_agent),
) -> TradingResponse:
    try:
        # Combine all data for the agent
        market_data = {
            "market": request.market_data.dict(),
            "sentiment": request.sentiment_data,
            "trend": request.trend_data,
        }
        
        # Get agent's decision
        result = agent.run(market_data)
        return TradingResponse(**result)
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Error analyzing trading opportunity: {str(e)}",
        )


@router.post("/execute")
async def execute_trade(
    request: TradingRequest,
    agent: TradingAgent = Depends(get_trading_agent),
) -> Dict[str, Any]:
    try:
        # Combine all data for the agent
        market_data = {
            "market": request.market_data.dict(),
            "sentiment": request.sentiment_data,
            "trend": request.trend_data,
        }
        
        # Get agent's decision and execute trade
        result = agent.run(market_data)
        if result["action"] == "buy" or result["action"] == "sell":
            trade_result = agent._execute_trade(result["parameters"])
            return {
                "status": "success",
                "trade": trade_result,
                "analysis": result,
            }
        return {
            "status": "no_action",
            "analysis": result,
        }
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Error executing trade: {str(e)}",
        ) 