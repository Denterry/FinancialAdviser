import os
import openai
from typing import List
from dotenv import load_dotenv

# Load environment variables from .env if present
load_dotenv()

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
MODEL_NAME = os.getenv("LLM_MODEL_NAME", "gpt-3.5-turbo")

if not OPENAI_API_KEY:
    raise RuntimeError("OPENAI_API_KEY not set in environment")

openai.api_key = OPENAI_API_KEY


async def call_openai_chat(messages: List[dict]) -> str:
    """
    Call OpenAI's chat completion endpoint and return the response.

    :param messages: List of messages, each with 'role' and 'content'
    :return: Generated string from the model
    """
    try:
        response = await openai.ChatCompletion.acreate(
            model=MODEL_NAME,
            messages=messages,
            temperature=0.7,
            max_tokens=512,
        )
        return response["choices"][0]["message"]["content"].strip()
    except openai.OpenAIError as e:
        raise RuntimeError(f"OpenAI API error: {e}")
