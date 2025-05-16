from typing import List, Dict, Any
import pandas as pd
import numpy as np
from transformers import AutoTokenizer, AutoModelForSequenceClassification
from prophet import Prophet
from langchain.agents import Tool, AgentExecutor, LLMSingleActionAgent
from langchain.prompts import StringPromptTemplate
from langchain import OpenAI, LLMChain
from langchain.tools import BaseTool

from src.core.config import settings


class SentimentAnalyzer:
    def __init__(self):
        self.tokenizer = AutoTokenizer.from_pretrained(settings.SENTIMENT_MODEL)
        self.model = AutoModelForSequenceClassification.from_pretrained(
            settings.SENTIMENT_MODEL
        )

    def analyze(self, text: str) -> Dict[str, float]:
        inputs = self.tokenizer(text, return_tensors="pt", truncation=True)
        outputs = self.model(**inputs)
        scores = outputs.logits.softmax(dim=1).detach().numpy()[0]
        return {
            "positive": float(scores[2]),
            "neutral": float(scores[1]),
            "negative": float(scores[0]),
        }


class TrendPredictor:
    def __init__(self):
        self.model = Prophet(
            yearly_seasonality=True,
            weekly_seasonality=True,
            daily_seasonality=True,
        )

    def predict(
        self, data: pd.DataFrame, periods: int = 30
    ) -> Dict[str, List[float]]:
        df = data.rename(columns={"date": "ds", "price": "y"})
        self.model.fit(df)
        future = self.model.make_future_dataframe(periods=periods)
        forecast = self.model.predict(future)
        return {
            "dates": forecast["ds"].tail(periods).dt.strftime("%Y-%m-%d").tolist(),
            "predictions": forecast["yhat"].tail(periods).tolist(),
            "lower_bound": forecast["yhat_lower"].tail(periods).tolist(),
            "upper_bound": forecast["yhat_upper"].tail(periods).tolist(),
        }


class TradingAgent:
    def __init__(self):
        self.llm = OpenAI(
            temperature=0,
            openai_api_key=settings.OPENAI_API_KEY,
            model_name=settings.OPENAI_MODEL,
        )
        self.tools = self._create_tools()
        self.agent = self._create_agent()

    def _create_tools(self) -> List[BaseTool]:
        return [
            Tool(
                name="GetMarketData",
                func=self._get_market_data,
                description="Get current market data for a symbol",
            ),
            Tool(
                name="AnalyzeSentiment",
                func=self._analyze_sentiment,
                description="Analyze sentiment of news and social media",
            ),
            Tool(
                name="PredictTrend",
                func=self._predict_trend,
                description="Predict price trend for a symbol",
            ),
            Tool(
                name="ExecuteTrade",
                func=self._execute_trade,
                description="Execute a trade with given parameters",
            ),
        ]

    def _create_agent(self) -> AgentExecutor:
        prompt = StringPromptTemplate.from_template(
            """You are an AI trading agent. Your goal is to make profitable trades based on market data, sentiment analysis, and trend predictions.

            Current market data: {market_data}
            Sentiment analysis: {sentiment}
            Trend prediction: {trend}

            What action should you take? Use the available tools to make a decision.

            {tools}

            {agent_scratchpad}"""
        )

        llm_chain = LLMChain(llm=self.llm, prompt=prompt)
        agent = LLMSingleActionAgent(
            llm_chain=llm_chain,
            allowed_tools=[tool.name for tool in self.tools],
            stop=["\nObservation:"],
        )

        return AgentExecutor.from_agent_and_tools(
            agent=agent,
            tools=self.tools,
            verbose=True,
        )

    def _get_market_data(self, symbol: str) -> Dict[str, Any]:
        # Implement market data retrieval
        pass

    def _analyze_sentiment(self, text: str) -> Dict[str, float]:
        analyzer = SentimentAnalyzer()
        return analyzer.analyze(text)

    def _predict_trend(self, symbol: str) -> Dict[str, List[float]]:
        # Implement trend prediction
        pass

    def _execute_trade(self, params: Dict[str, Any]) -> Dict[str, Any]:
        if not settings.TRADING_ENABLED:
            raise ValueError("Trading is not enabled")
        # Implement trade execution
        pass

    def run(self, market_data: Dict[str, Any]) -> Dict[str, Any]:
        return self.agent.run(market_data) 