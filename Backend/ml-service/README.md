# ML Service

A Python-based microservice for financial analysis and automated trading using machine learning and LLMs.

## Features

- **Sentiment Analysis**: Analyze market sentiment using FinBERT model
- **Trend Prediction**: Predict price trends using Prophet
- **AI Trading Agent**: Automated trading using LLMs and various tools
- **gRPC Integration**: Connect with other services
- **REST API**: FastAPI endpoints for ML operations

## Prerequisites

- Python 3.9+
- PostgreSQL
- OpenAI API key
- Trading API credentials (optional)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/FinancialAdviser.git
cd FinancialAdviser/Backend/ml-service
```

2. Create and activate virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # Linux/Mac
# or
.\venv\Scripts\activate  # Windows
```

3. Install dependencies:
```bash
pip install -r requirements.txt
```

4. Copy environment file and configure:
```bash
cp .env.example .env
# Edit .env with your configuration
```

## Usage

1. Start the service:
```bash
uvicorn src.app:app --reload
```

2. Access the API documentation:
- Swagger UI: http://localhost:8000/api/docs
- ReDoc: http://localhost:8000/api/redoc

## API Endpoints

### Sentiment Analysis
- `POST /api/v1/sentiment/analyze`: Analyze text sentiment
- `POST /api/v1/sentiment/batch`: Batch sentiment analysis

### Trend Prediction
- `POST /api/v1/trends/predict`: Predict price trends
- `GET /api/v1/trends/symbols`: Get available symbols

### Trading
- `POST /api/v1/trading/analyze`: Analyze trading opportunities
- `POST /api/v1/trading/execute`: Execute trades

## ML Models

### Sentiment Analysis
- Uses FinBERT model for financial sentiment analysis
- Provides positive, neutral, and negative scores

### Trend Prediction
- Uses Prophet for time series forecasting
- Supports multiple symbols and timeframes

### Trading Agent
- Powered by GPT-4
- Uses various tools for market analysis
- Makes trading decisions based on multiple factors

## Development

### Running Tests
```bash
pytest
```

### Code Style
```bash
black .
isort .
flake8
mypy .
```

## Docker

Build and run with Docker:
```bash
docker-compose up --build
```

## License

MIT
