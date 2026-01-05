from fastapi import APIRouter
import structlog

from app.domain.schemas import (
    RecommendationRequest,
    RecommendationResponse,
    EmbeddingRequest,
    EmbeddingResponse,
)
from app.ml.embeddings import embedding_service
from app.ml.recommendations import recommendation_engine

logger = structlog.get_logger()
router = APIRouter()


@router.post("/recommendations", response_model=list[RecommendationResponse])
async def get_recommendations(request: RecommendationRequest):
    logger.info(
        "Recommendation request",
        user_id=str(request.user_id),
        recommendation_type=request.recommendation_type,
    )

    recommendations = []

    if request.recommendation_type == "content_based":
        recs = recommendation_engine.get_content_based_recommendations(
            query_text="",
            limit=request.limit,
        )
    elif request.recommendation_type == "collaborative":
        recs = recommendation_engine.get_collaborative_recommendations(
            user_id=str(request.user_id),
            limit=request.limit,
        )
    else:
        recs = recommendation_engine.get_hybrid_recommendations(
            user_id=str(request.user_id),
            limit=request.limit,
        )

    for rec in recs:
        recommendations.append(
            RecommendationResponse(
                source_id=rec["source_id"],
                title=rec["metadata"].get("title", "Unknown"),
                score=rec["score"],
                reason="Based on your reading history",
            )
        )

    return recommendations


@router.post("/embeddings", response_model=EmbeddingResponse)
async def create_embedding(request: EmbeddingRequest):
    logger.info("Embedding request", text_length=len(request.text))

    vector = embedding_service.embed(request.text)

    return EmbeddingResponse(
        vector=vector,
        dimension=len(vector),
    )
