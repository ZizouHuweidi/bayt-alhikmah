import uuid
from datetime import datetime
import structlog

from fastapi import APIRouter, HTTPException, Depends

from app.domain.schemas import ScrapeRequest, ScrapeResponse
from app.scraper.base import URLScraper
from app.core.kafka import KafkaProducerManager
from app.core.config import settings

logger = structlog.get_logger()
router = APIRouter()


async def get_kafka_producer() -> KafkaProducerManager:
    from main import kafka_producer

    return kafka_producer


@router.post("/scrape", response_model=ScrapeResponse)
async def scrape_url(
    request: ScrapeRequest,
    kafka_producer: KafkaProducerManager = Depends(get_kafka_producer),
):
    event_id = str(uuid.uuid4())

    logger.info("Received scrape request", event_id=event_id, url=request.url)

    scraper = URLScraper(
        user_agent=settings.scraper_user_agent,
        timeout=settings.scraper_timeout,
    )

    try:
        metadata = await scraper.scrape_url(request.url)
        await scraper.close()

        event = {
            "event_id": event_id,
            "url": request.url,
            "source_type": metadata.get("source_type"),
            "scraped_at": datetime.utcnow().isoformat(),
            "raw_html": None,
            "metadata": metadata,
        }

        import json

        await kafka_producer.publish(
            topic=settings.kafka_topic_scrape_raw,
            value=json.dumps(event).encode("utf-8"),
            key=event_id.encode("utf-8"),
        )

        return ScrapeResponse(
            success=True,
            message="Scraping completed successfully",
            event_id=event_id,
        )

    except Exception as e:
        logger.error("Scraping failed", event_id=event_id, error=str(e))
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/scrape/isbn/{isbn}")
async def scrape_isbn(
    isbn: str,
    kafka_producer: KafkaProducerManager = Depends(get_kafka_producer),
):
    event_id = str(uuid.uuid4())

    logger.info("Received ISBN scrape request", event_id=event_id, isbn=isbn)

    scraper = URLScraper(user_agent=settings.scraper_user_agent)

    try:
        metadata = await scraper.scrape_isbn(isbn)
        await scraper.close()

        event = {
            "event_id": event_id,
            "url": f"isbn:{isbn}",
            "source_type": metadata.get("source_type"),
            "scraped_at": datetime.utcnow().isoformat(),
            "raw_html": None,
            "metadata": metadata,
        }

        import json

        await kafka_producer.publish(
            topic=settings.kafka_topic_scrape_raw,
            value=json.dumps(event).encode("utf-8"),
            key=event_id.encode("utf-8"),
        )

        return ScrapeResponse(
            success=True,
            message=f"ISBN {isbn} scraped successfully",
            event_id=event_id,
        )

    except Exception as e:
        logger.error("ISBN scraping failed", event_id=event_id, isbn=isbn, error=str(e))
        raise HTTPException(status_code=500, detail=str(e))
