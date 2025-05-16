import threading
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.core.config import settings
from src.api.v1.router import api_router
from src.core.logging import setup_logging
from src.grpc.server import serve as serve_grpc


app = FastAPI(
    title=settings.PROJECT_NAME,
    version=settings.VERSION,
    description="ML Service for Financial Analysis and Trading",
    docs_url="/api/docs",
    redoc_url="/api/redoc",
    openapi_url="/api/openapi.json",
)

# Setup CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include API router
app.include_router(api_router, prefix="/api/v1")

# Setup logging
setup_logging()


@app.get("/health")
async def health_check():
    return {"status": "healthy"}


def start_grpc_server():
    """Start the gRPC server in a separate thread."""
    grpc_thread = threading.Thread(target=serve_grpc)
    grpc_thread.daemon = True
    grpc_thread.start()


@app.on_event("startup")
async def startup_event():
    """Start the gRPC server when the FastAPI application starts."""
    start_grpc_server()
