from typing import List
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
import pandas as pd

from src.ml.models import TrendPredictor

router = APIRouter()
trend_predictor = TrendPredictor()


class PriceData(BaseModel):
    date: str
    price: float


class TrendRequest(BaseModel):
    symbol: str
    data: List[PriceData]
    periods: int = 30


class TrendResponse(BaseModel):
    dates: List[str]
    predictions: List[float]
    lower_bound: List[float]
    upper_bound: List[float]


@router.post("/predict", response_model=TrendResponse)
async def predict_trend(request: TrendRequest) -> TrendResponse:
    try:
        # Convert request data to DataFrame
        df = pd.DataFrame([d.dict() for d in request.data])
        
        # Make prediction
        result = trend_predictor.predict(df, request.periods)
        return TrendResponse(**result)
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Error predicting trend: {str(e)}",
        )


@router.get("/symbols")
async def get_available_symbols() -> List[str]:
    # Implement symbol list retrieval
    return ["AAPL", "GOOGL", "MSFT", "AMZN", "META"] 