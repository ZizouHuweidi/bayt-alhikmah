import sys
from contextlib import asynccontextmanager

import structlog
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.api import scrape_router
from app.core.config import settings
from app.core.logging import configure_logging
from app.core.kafka import KafkaProducerManager

configure_logging()
logger = structlog.get_logger()

kafka_producer = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global kafka_producer
    logger.info("Starting Bahith scraper service")

    kafka_producer = KafkaProducerManager(
        bootstrap_servers=settings.kafka_brokers,
    )
    await kafka_producer.start()

    logger.info("Bahith service started")
    yield

    logger.info("Shutting down Bahith service")
    await kafka_producer.stop()
    logger.info("Bahith service stopped")


app = FastAPI(
    title="Bahith Scraper Service",
    description="Scraper and ingestion service for Bayt al Hikmah",
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

app.include_router(scrape_router, prefix="/api/v1", tags=["scrape"])


@app.get("/healthz")
async def liveness():
    return {"status": "alive"}


@app.get("/readyz")
async def readiness():
    global kafka_producer

    if kafka_producer is None or not kafka_producer.is_connected():
        return {"status": "not ready", "error": "Kafka not connected"}

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
        port=8003,
        reload=True,
        log_level="info",
    )
