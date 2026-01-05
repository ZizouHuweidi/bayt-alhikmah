from typing import List, Optional
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity
import structlog

from app.ml.embeddings import embedding_service
from app.domain.schemas import RecommendationType

logger = structlog.get_logger()


class RecommendationEngine:
    def __init__(self):
        self.source_embeddings: dict = {}
        self.source_metadata: dict = {}

    def add_source(self, source_id: str, text: str, metadata: dict):
        embedding = embedding_service.embed(text)
        self.source_embeddings[source_id] = embedding
        self.source_metadata[source_id] = metadata
        logger.info("Source added to recommendation engine", source_id=source_id)

    def get_content_based_recommendations(
        self,
        query_text: str,
        limit: int = 10,
    ) -> List[dict]:
        if not self.source_embeddings:
            return []

        query_embedding = embedding_service.embed(query_text)
        source_ids = list(self.source_embeddings.keys())
        source_vectors = np.array([self.source_embeddings[sid] for sid in source_ids])

        similarities = cosine_similarity(
            [query_embedding],
            source_vectors,
        )[0]

        top_indices = np.argsort(similarities)[::-1][:limit]

        recommendations = []
        for idx in top_indices:
            source_id = source_ids[idx]
            recommendations.append(
                {
                    "source_id": source_id,
                    "score": float(similarities[idx]),
                    "metadata": self.source_metadata[source_id],
                }
            )

        return recommendations

    def get_collaborative_recommendations(
        self,
        user_id: str,
        limit: int = 10,
    ) -> List[dict]:
        logger.info("Collaborative filtering not implemented yet", user_id=user_id)
        return []

    def get_hybrid_recommendations(
        self,
        user_id: str,
        query_text: Optional[str] = None,
        limit: int = 10,
    ) -> List[dict]:
        content_based = []
        if query_text:
            content_based = self.get_content_based_recommendations(query_text, limit)

        collaborative = self.get_collaborative_recommendations(user_id, limit)

        scores = {}
        for rec in content_based:
            source_id = rec["source_id"]
            scores[source_id] = scores.get(source_id, 0) + rec["score"] * 0.7

        for rec in collaborative:
            source_id = rec["source_id"]
            scores[source_id] = scores.get(source_id, 0) + rec["score"] * 0.3

        sorted_recs = sorted(scores.items(), key=lambda x: x[1], reverse=True)[:limit]

        return [
            {
                "source_id": source_id,
                "score": score,
                "metadata": self.source_metadata.get(source_id, {}),
            }
            for source_id, score in sorted_recs
        ]


recommendation_engine = RecommendationEngine()
