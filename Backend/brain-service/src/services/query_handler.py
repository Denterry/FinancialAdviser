from typing import Dict, List
import httpx
from app.core.config import settings
from app.core.llm_client import llm_client
from app.services.prompt_builder import PromptBuilder
from app.utils.ticker_extractor import extract_tickers


class QueryHandler:
    def __init__(self):
        self.prompt_builder = PromptBuilder()
        self.client = httpx.AsyncClient()

    async def handle_query(
        self,
        user_id: str,
        query: str,
        user_context: Dict = None
    ) -> Dict:
        """
        Handle a user query by processing it and generating a response.
        
        Args:
            user_id: The ID of the user making the query
            query: The user's query text
            user_context: Optional user context (risk profile, preferences)
            
        Returns:
            A dictionary containing the analysis and recommendations
        """
        try:
            # Extract tickers from the query
            tickers = extract_tickers(query)
            if not tickers:
                return {
                    "error": (
                        "No financial instruments found in the query. "
                        "Please mention specific stocks or assets."
                    )
                }

            # Fetch market data for each ticker
            market_data = await self._fetch_market_data(tickers)

            # Build and get analysis
            analysis_prompt = self.prompt_builder.build_analysis_prompt(
                query, market_data, user_context
            )
            analysis = await llm_client.get_completion(analysis_prompt)

            # Build and get recommendations
            if user_context:
                rec_prompt = self.prompt_builder.build_recommendation_prompt(
                    query, market_data, user_context
                )
                recommendations = await llm_client.get_completion(rec_prompt)
            else:
                recommendations = (
                    "Please provide your risk profile and investment "
                    "preferences to receive personalized recommendations."
                )

            return {
                "analysis": analysis,
                "recommendations": recommendations,
                "market_data": market_data
            }

        except Exception as e:
            return {"error": str(e)}
        finally:
            await self.client.aclose()

    async def _fetch_market_data(self, tickers: List[str]) -> Dict[str, Dict]:
        """
        Fetch market data for the given tickers from external services.
        
        Args:
            tickers: List of ticker symbols
            
        Returns:
            Dictionary containing market data for each ticker
        """
        market_data = {}
        
        for ticker in tickers:
            # Fetch sentiment data
            sentiment_response = await self.client.get(
                f"{settings.SENTIMENT_SERVICE_URL}/sentiment/{ticker}"
            )
            sentiment_data = sentiment_response.json()

            # Fetch price prediction
            prediction_response = await self.client.get(
                f"{settings.PREDICTION_SERVICE_URL}/predict/{ticker}"
            )
            prediction_data = prediction_response.json()

            market_data[ticker] = {
                "current_price": prediction_data.get("current_price"),
                "price_change_24h": prediction_data.get("price_change_24h"),
                "sentiment": sentiment_data.get("sentiment"),
                "prediction": prediction_data.get("prediction")
            }

        return market_data 