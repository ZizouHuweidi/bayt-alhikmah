from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    server_port: int = 8003
    environment: str = "development"

    kafka_brokers: str = "redpanda:9092"
    kafka_topic_scrape_raw: str = "scrape.raw"
    kafka_topic_scrape_processed: str = "scrape.processed"

    scraper_timeout: int = 30
    scraper_user_agent: str = "Bayt-al-Hikmah-Bahith/1.0"
    scraper_max_retries: int = 3

    database_url: str = "sqlite:///./bahith.db"

    otel_exporter_endpoint: str = "http://tempo:4317"
    otel_service_name: str = "bayt-alhikmah-bahith"

    log_level: str = "info"

    class Config:
        env_file = ".env"
        case_sensitive = False


@lru_cache()
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
