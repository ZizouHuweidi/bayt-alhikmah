import asyncio
from typing import Optional
from aiokafka import AIOKafkaProducer
import structlog

from app.core.config import settings

logger = structlog.get_logger()


class KafkaProducerManager:
    def __init__(self, bootstrap_servers: str):
        self.bootstrap_servers = bootstrap_servers
        self.producer: Optional[AIOKafkaProducer] = None
        self._connected = False

    async def start(self):
        self.producer = AIOKafkaProducer(
            bootstrap_servers=self.bootstrap_servers,
        )
        await self.producer.start()
        self._connected = True
        logger.info("Kafka producer started", bootstrap_servers=self.bootstrap_servers)

    async def stop(self):
        if self.producer:
            await self.producer.stop()
            self._connected = False
            logger.info("Kafka producer stopped")

    def is_connected(self) -> bool:
        return self._connected

    async def publish(self, topic: str, value: bytes, key: Optional[bytes] = None):
        if not self._connected:
            logger.error("Cannot publish: Kafka producer not connected")
            return False

        try:
            await self.producer.send_and_wait(topic, value=value, key=key)
            logger.info("Message published", topic=topic, key=key.decode() if key else None)
            return True
        except Exception as e:
            logger.error("Failed to publish message", topic=topic, error=str(e))
            return False
