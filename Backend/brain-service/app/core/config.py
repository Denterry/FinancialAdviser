from typing import List
from pydantic_settings import BaseSettings
from pydantic import AnyHttpUrl

class Settings(BaseSettings):
    # API settings
    API_V1_STR: str = "/api/v1"
    PROJECT_NAME: str = "Brain Service"
    
    # CORS settings
    CORS_ORIGINS: List[AnyHttpUrl] = ["http://localhost:3000"]  # Frontend URL
    
    # LLM settings
    OPENAI_API_KEY: str
    OPENAI_MODEL: str = "gpt-3.5-turbo"
    
    # External services
    SENTIMENT_SERVICE_URL: str = "http://sentiment-service:50051"
    PREDICTION_SERVICE_URL: str = "http://prediction-service:50051"
    
    # Database settings
    POSTGRES_HOST: str = "localhost"
    POSTGRES_PORT: int = 5432
    POSTGRES_USER: str = "postgres"
    POSTGRES_PASSWORD: str = "postgres"
    POSTGRES_DB: str = "brain_service"
    
    class Config:
        case_sensitive = True
        env_file = ".env"

settings = Settings() 