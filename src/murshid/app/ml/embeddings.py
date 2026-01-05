from typing import List
import numpy as np
from sentence_transformers import SentenceTransformer
import structlog

from app.core.config import settings

logger = structlog.get_logger()


class EmbeddingService:
    def __init__(self, model_name: str, dimension: int):
        self.model_name = model_name
        self.dimension = dimension
        self.model: SentenceTransformer = None

    def load_model(self):
        logger.info("Loading embedding model", model=self.model_name)
        self.model = SentenceTransformer(self.model_name)
        logger.info("Embedding model loaded")

    def embed(self, text: str) -> List[float]:
        if self.model is None:
            self.load_model()

        embedding = self.model.encode(text, convert_to_numpy=True)
        return embedding.tolist()

    def embed_batch(self, texts: List[str]) -> List[List[float]]:
        if self.model is None:
            self.load_model()

        embeddings = self.model.encode(texts, convert_to_numpy=True)
        return [e.tolist() for e in embeddings]


embedding_service = EmbeddingService(
    model_name=settings.embeddings_model,
    dimension=settings.vector_dimension,
)
