from pydantic import BaseModel, HttpUrl
from typing import Optional, List
from datetime import datetime
from enum import Enum
import uuid


class RecommendationType(str, Enum):
    CONTENT_BASED = "content_based"
    COLLABORATIVE = "collaborative"
    HYBRID = "hybrid"


class RecommendationRequest(BaseModel):
    user_id: uuid.UUID
    recommendation_type: RecommendationType = RecommendationType.HYBRID
    limit: int = 10


class RecommendationResponse(BaseModel):
    source_id: uuid.UUID
    title: str
    score: float
    reason: str


class EmbeddingRequest(BaseModel):
    text: str


class EmbeddingResponse(BaseModel):
    vector: List[float]
    dimension: int
