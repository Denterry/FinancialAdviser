from typing import List
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from src.ml.models import SentimentAnalyzer

router = APIRouter()
sentiment_analyzer = SentimentAnalyzer()


class SentimentRequest(BaseModel):
    text: str


class SentimentResponse(BaseModel):
    positive: float
    neutral: float
    negative: float


@router.post("/analyze", response_model=SentimentResponse)
async def analyze_sentiment(request: SentimentRequest) -> SentimentResponse:
    try:
        result = sentiment_analyzer.analyze(request.text)
        return SentimentResponse(**result)
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Error analyzing sentiment: {str(e)}",
        )


@router.post("/batch", response_model=List[SentimentResponse])
async def analyze_sentiment_batch(
    requests: List[SentimentRequest]
) -> List[SentimentResponse]:
    try:
        results = []
        for request in requests:
            result = sentiment_analyzer.analyze(request.text)
            results.append(SentimentResponse(**result))
        return results
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Error analyzing sentiment batch: {str(e)}",
        ) 