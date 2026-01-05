from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    server_port: int = 8004
    environment: str = "development"

    kafka_brokers: str = "redpanda:9092"
    kafka_topic_source_created: str = "source.created"
    kafka_topic_note_created: str = "note.created"

    database_url: str = (
        "postgres://postgres:postgres_password_change_in_production@postgres:5432/murshid?sslmode=disable"
    )

    embeddings_model: str = "sentence-transformers/all-MiniLM-L6-v2"
    vector_dimension: int = 384
    recommendations_min_interactions: int = 5

    otel_exporter_endpoint: str = "http://tempo:4317"
    otel_service_name: str = "bayt-alhikmah-murshid"

    log_level: str = "info"

    class Config:
        env_file = ".env"
        case_sensitive = False


@lru_cache()
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
