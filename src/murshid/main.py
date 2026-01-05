import sys
from contextlib import asynccontextmanager

import structlog
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.api import recommendations_router
from app.core.config import settings
from app.core.logging import configure_logging

configure_logging()
logger = structlog.get_logger()


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("Starting Murshid ML service")

    logger.info("Murshid service started")
    yield

    logger.info("Shutting down Murshid service")
    logger.info("Murshid service stopped")


app = FastAPI(
    title="Murshid ML Service",
    description="ML and recommendations service for Bayt al Hikmah",
    version="0.1.0",
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(recommendations_router, prefix="/api/v1", tags=["recommendations"])


@app.get("/healthz")
async def liveness():
    return {"status": "alive"}


@app.get("/readyz")
async def readiness():
    return {"status": "ready"}


@app.get("/metrics")
async def metrics():
    from prometheus_client import generate_latest, CONTENT_TYPE_LATEST
    from fastapi import Response

    return Response(content=generate_latest(), media_type=CONTENT_TYPE_LATEST)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8004,
        reload=True,
        log_level="info",
    )
