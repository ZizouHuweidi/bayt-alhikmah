from pydantic import BaseModel, HttpUrl
from typing import Optional, List
from datetime import datetime
from enum import Enum


class SourceType(str, Enum):
    BOOK = "book"
    PAPER = "paper"
    PODCAST = "podcast"
    VIDEO = "video"
    ARTICLE = "article"
    ESSAY = "essay"


class ScrapeRequest(BaseModel):
    url: str
    source_type: Optional[SourceType] = None
    force_refresh: bool = False


class ScrapeRawEvent(BaseModel):
    event_id: str
    url: str
    source_type: Optional[SourceType] = None
    scraped_at: datetime
    raw_html: Optional[str] = None
    metadata: Optional[dict] = None


class ScrapeProcessedEvent(BaseModel):
    event_id: str
    url: str
    source_type: SourceType
    title: str
    description: Optional[str] = None
    author: Optional[str] = None
    publisher: Optional[str] = None
    isbn: Optional[str] = None
    doi: Optional[str] = None
    tags: List[str] = []
    published_at: Optional[datetime] = None
    processed_at: datetime


class ScrapeResponse(BaseModel):
    success: bool
    message: str
    event_id: Optional[str] = None
