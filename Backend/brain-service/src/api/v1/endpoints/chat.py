from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import Optional

from services.chat_orchestrator import ChatOrchestrator

router = APIRouter()


class ChatRequest(BaseModel):
    user_id: str
    message: str
    chat_id: Optional[int] = None


class ChatResponse(BaseModel):
    chat_id: int
    response: str


@router.post("/chat", response_model=ChatResponse)
async def handle_chat(req: ChatRequest):
    try:
        orchestrator = ChatOrchestrator()
        response, chat_id = await orchestrator.process_user_query(
            user_id=req.user_id,
            message=req.message,
            chat_id=req.chat_id
        )
        return ChatResponse(chat_id=chat_id, response=response)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
