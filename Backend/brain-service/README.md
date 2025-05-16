# Financial Adviser Brain Service

AI-powered financial analysis and recommendations service that processes user queries about financial instruments and provides intelligent insights and recommendations.

## Features

- Natural language processing of financial queries
- Integration with external services for market data and sentiment analysis
- AI-powered analysis and recommendations using OpenAI's GPT models
- RESTful API for easy integration with other services

## Prerequisites

- Python 3.8+
- OpenAI API key
- Access to external services (sentiment analysis and price prediction)

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd brain-service
```

2. Create and activate a virtual environment:
```bash
python -m venv venv
source venv/bin/activate
```

3. Install dependencies:
```bash
pip install -r requirements.txt
```

4. Create a `.env` file in the root directory with the following variables:
```env
OPENAI_API_KEY=your_api_key_here
SENTIMENT_SERVICE_URL=http://sentiment-service:8000
PREDICTION_SERVICE_URL=http://prediction-service:8000
CORS_ORIGINS=["http://localhost:3000"]
```

## Running the Service

Start the service with:
```bash
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

The service will be available at `http://localhost:8000`.

## API Documentation

Once the service is running, you can access the interactive API documentation at:
- Swagger UI: `http://localhost:8000/docs`
- ReDoc: `http://localhost:8000/redoc`

### Main Endpoints

- `POST /api/v1/query`: Process a financial query and get analysis/recommendations
  ```json
  {
    "user_id": "string",
    "query": "string",
    "user_context": {
      "risk_profile": "string",
      "investment_preferences": ["string"]
    }
  }
  ```

- `GET /health`: Health check endpoint

## Development

### Adding New Features

1. Create new modules in the appropriate directories
2. Update the router to include new endpoints
3. Add new dependencies to requirements.txt if needed
4. Update the README with new features and documentation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request


### Project Structure

```
brain-service/
├── docker-compose-integration-test.yml
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── migrations
│   ├── 20250525134008_uuid_ext.down.sql
│   ├── 20250525134008_uuid_ext.up.sql
│   ├── 20250525134009_chats_table.down.sql
│   ├── 20250525134009_chats_table.up.sql
│   ├── 20250525134026_messages_table.down.sql
│   ├── 20250525134026_messages_table.up.sql
│   ├── 20250609180124_update_updated_at_trigger.down.sql
│   ├── 20250609180124_update_updated_at_trigger.up.sql
│   ├── 20250609181040_llm_logs_table.down.sql
│   ├── 20250609181040_llm_logs_table.up.sql
│   ├── 20250609181542_user_usage_quota_table.down.sql
│   └── 20250609181542_user_usage_quota_table.up.sql
├── README.md
├── requirements.txt
└── src
    ├── core
    │   └── db.py
    ├── grpc
    │   └── server.py
    ├── main.py
    ├── proto
    │   └── brain
    │       └── brain.proto
    ├── repo
    │   ├── chats.py
    │   └── messages.py
    ├── services
    │   ├── llm_orchestrator.py
    │   └── ml_client.py
    └── utils
        ├── prompt_builder.py
        ├── ticker_extractor.py
        └── tokenizer.py
```
