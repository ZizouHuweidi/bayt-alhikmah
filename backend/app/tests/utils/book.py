from sqlmodel import Session

from app import crud
from app.models import Book, BookCreate
from app.tests.utils.user import create_random_user
from app.tests.utils.utils import random_lower_string


def create_random_book(db: Session) -> Book:
    user = create_random_user(db)
    owner_id = user.id
    assert owner_id is not None
    title = random_lower_string()
    description = random_lower_string()
    book_in = BookCreate(title=title, description=description)
    return crud.create_book(session=db, book_in=book_in, owner_id=owner_id)
