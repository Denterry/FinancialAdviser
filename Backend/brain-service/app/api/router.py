from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import Dict, Optional
from app.services.query_handler import QueryHandler

router = APIRouter()
query_handler = QueryHandler()


class QueryRequest(BaseModel):
    user_id: str
    query: str
    user_context: Optional[Dict] = None


class QueryResponse(BaseModel):
    analysis: str
    recommendations: str
    market_data: Dict


@router.post("/query", response_model=QueryResponse)
async def handle_query(request: QueryRequest):
    """
    Handle a user query and return analysis and recommendations.
    
    Args:
        request: The query request containing user ID, query text, and optional
                context
        
    Returns:
        Analysis, recommendations, and market data for the query
    """
    try:
        result = await query_handler.handle_query(
            request.user_id,
            request.query,
            request.user_context
        )
        
        if "error" in result:
            raise HTTPException(status_code=400, detail=result["error"])
            
        return result
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) 