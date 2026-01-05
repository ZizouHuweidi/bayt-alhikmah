import httpx
import structlog
from typing import Optional, Dict, Any
from bs4 import BeautifulSoup

from app.domain.schemas import SourceType

logger = structlog.get_logger()


class BaseScraper:
    def __init__(self, user_agent: str, timeout: int = 30):
        self.user_agent = user_agent
        self.timeout = timeout
        self.client = httpx.AsyncClient(
            headers={"User-Agent": user_agent},
            timeout=timeout,
        )

    async def close(self):
        await self.client.aclose()

    async def fetch_url(self, url: str) -> str:
        try:
            response = await self.client.get(url)
            response.raise_for_status()
            return response.text
        except Exception as e:
            logger.error("Failed to fetch URL", url=url, error=str(e))
            raise

    def extract_metadata(self, html: str) -> Dict[str, Any]:
        soup = BeautifulSoup(html, "html.parser")

        metadata = {
            "title": None,
            "description": None,
            "author": None,
        }

        title_tag = soup.find("title")
        if title_tag:
            metadata["title"] = title_tag.get_text().strip()

        meta_description = soup.find("meta", attrs={"name": "description"})
        if meta_description:
            metadata["description"] = meta_description.get("content")

        meta_author = soup.find("meta", attrs={"name": "author"})
        if meta_author:
            metadata["author"] = meta_author.get("content")

        return metadata


class BookScraper(BaseScraper):
    async def scrape_isbn(self, isbn: str) -> Dict[str, Any]:
        logger.info("Scraping book by ISBN", isbn=isbn)

        html = await self.fetch_url(f"https://openlibrary.org/isbn/{isbn}")
        metadata = self.extract_metadata(html)

        metadata["isbn"] = isbn
        metadata["source_type"] = SourceType.BOOK

        return metadata


class URLScraper(BaseScraper):
    async def scrape_url(self, url: str) -> Dict[str, Any]:
        logger.info("Scraping URL", url=url)

        html = await self.fetch_url(url)
        metadata = self.extract_metadata(html)

        metadata["url"] = url
        metadata["source_type"] = self._guess_source_type(url)

        return metadata

    def _guess_source_type(self, url: str) -> SourceType:
        url_lower = url.lower()

        if "arxiv.org" in url_lower:
            return SourceType.PAPER
        elif "youtube.com" in url_lower or "youtu.be" in url_lower:
            return SourceType.VIDEO
        elif any(domain in url_lower for domain in ["podcasts.", "podcast."]):
            return SourceType.PODCAST
        else:
            return SourceType.ARTICLE
