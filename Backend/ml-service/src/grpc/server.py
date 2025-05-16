import grpc
from concurrent import futures
import logging

from src.ml.models import SentimentAnalyzer, TrendPredictor, TradingAgent
from src.core.config import settings
from src.proto import ml_service_pb2, ml_service_pb2_grpc


class MLServiceServicer(ml_service_pb2_grpc.MLServiceServicer):
    def __init__(self):
        self.sentiment_analyzer = SentimentAnalyzer()
        self.trend_predictor = TrendPredictor()
        self.trading_agent = TradingAgent()

    def AnalyzeSentiment(self, request, context):
        try:
            result = self.sentiment_analyzer.analyze(request.text)
            return ml_service_pb2.SentimentResponse(
                positive=result["positive"],
                neutral=result["neutral"],
                negative=result["negative"],
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return ml_service_pb2.SentimentResponse()

    def BatchAnalyzeSentiment(self, request, context):
        try:
            results = []
            for text in request.texts:
                result = self.sentiment_analyzer.analyze(text)
                results.append(
                    ml_service_pb2.SentimentResponse(
                        positive=result["positive"],
                        neutral=result["neutral"],
                        negative=result["negative"],
                    )
                )
            return ml_service_pb2.BatchSentimentResponse(results=results)
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return ml_service_pb2.BatchSentimentResponse()

    def PredictTrend(self, request, context):
        try:
            # Convert request data to DataFrame
            import pandas as pd
            data = {
                "date": [d.date for d in request.data],
                "price": [d.price for d in request.data],
            }
            df = pd.DataFrame(data)
            
            # Make prediction
            result = self.trend_predictor.predict(df, request.periods)
            return ml_service_pb2.TrendResponse(
                dates=result["dates"],
                predictions=result["predictions"],
                lower_bound=result["lower_bound"],
                upper_bound=result["upper_bound"],
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return ml_service_pb2.TrendResponse()

    def GetSymbols(self, request, context):
        try:
            symbols = ["AAPL", "GOOGL", "MSFT", "AMZN", "META"]
            return ml_service_pb2.GetSymbolsResponse(symbols=symbols)
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return ml_service_pb2.GetSymbolsResponse()

    def AnalyzeTrading(self, request, context):
        try:
            # Convert request data
            market_data = {
                "market": {
                    "symbol": request.market_data.symbol,
                    "price": request.market_data.price,
                    "volume": request.market_data.volume,
                    "timestamp": request.market_data.timestamp,
                    "indicators": dict(request.market_data.indicators),
                },
                "sentiment": dict(request.sentiment_data),
                "trend": dict(request.trend_data),
            }
            
            # Get agent's decision
            result = self.trading_agent.run(market_data)
            return ml_service_pb2.TradingResponse(
                action=result["action"],
                confidence=result["confidence"],
                parameters=result["parameters"],
                explanation=result["explanation"],
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return ml_service_pb2.TradingResponse()

    def ExecuteTrade(self, request, context):
        try:
            # Convert request data
            market_data = {
                "market": {
                    "symbol": request.market_data.symbol,
                    "price": request.market_data.price,
                    "volume": request.market_data.volume,
                    "timestamp": request.market_data.timestamp,
                    "indicators": dict(request.market_data.indicators),
                },
                "sentiment": dict(request.sentiment_data),
                "trend": dict(request.trend_data),
            }
            
            # Get agent's decision and execute trade
            result = self.trading_agent.run(market_data)
            if result["action"] in ["buy", "sell"]:
                trade_result = self.trading_agent._execute_trade(result["parameters"])
                return ml_service_pb2.TradeExecutionResponse(
                    status="success",
                    trade=trade_result,
                    analysis=ml_service_pb2.TradingResponse(
                        action=result["action"],
                        confidence=result["confidence"],
                        parameters=result["parameters"],
                        explanation=result["explanation"],
                    ),
                )
            return ml_service_pb2.TradeExecutionResponse(
                status="no_action",
                analysis=ml_service_pb2.TradingResponse(
                    action=result["action"],
                    confidence=result["confidence"],
                    parameters=result["parameters"],
                    explanation=result["explanation"],
                ),
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return ml_service_pb2.TradeExecutionResponse()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ml_service_pb2_grpc.add_MLServiceServicer_to_server(
        MLServiceServicer(), server
    )
    server.add_insecure_port(f"[::]:{settings.GRPC_PORT}")
    server.start()
    logging.info(f"gRPC server started on port {settings.GRPC_PORT}")
    server.wait_for_termination() 