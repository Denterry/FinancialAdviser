import openai
from app.core.config import settings


class LLMClient:
    def __init__(self):
        openai.api_key = settings.OPENAI_API_KEY
        self.model = settings.OPENAI_MODEL

    async def get_completion(self, prompt: str) -> str:
        """
        Get a completion from the LLM.
        
        Args:
            prompt: The prompt to send to the LLM
            
        Returns:
            The LLM's response
        """
        try:
            response = await openai.ChatCompletion.acreate(
                model=self.model,
                messages=[
                    {"role": "system", "content": "You are a financial advisor AI. Provide insightful, concise advice based on data."},
                    {"role": "user", "content": prompt}
                ],
                temperature=0.7,
                max_tokens=500
            )
            return response.choices[0].message["content"]
        except Exception as e:
            raise Exception(f"Error getting LLM completion: {str(e)}")


# Create a singleton instance
llm_client = LLMClient() 