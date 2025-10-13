from fastapi import APIRouter, HTTPException, Query
from sqlalchemy.orm import Session
from typing import List, Optional
from models import Book
from pydantic import BaseModel
from database import SessionLocal

router = APIRouter(prefix="/books", tags=["Books"])

class BookSchema(BaseModel):
    id: Optional[int]
    title: str
    genre: Optional[str] = None
    year: Optional[int] = None

    class Config:
        orm_mode = True


@router.get("/", response_model=List[BookSchema])
def list_books(
    title: Optional[str] = None,
    genre: Optional[str] = None,
    year: Optional[int] = None,
    limit: int = Query(10, ge=1),
    offset: int = Query(0, ge=0),
):
    with SessionLocal() as db:
        query = db.query(Book)
        if title:
            query = query.filter(Book.title.ilike(f"%{title}%"))
        if genre:
            query = query.filter(Book.genre.ilike(f"%{genre}%"))
        if year:
            query = query.filter(Book.year == year)
        return query.order_by(Book.id).limit(limit).offset(offset).all()


@router.post("/", response_model=BookSchema, status_code=201)
def create_book(book: BookSchema):
    with SessionLocal() as db:
        new_book = Book(title=book.title, genre=book.genre, year=book.year)
        db.add(new_book)
        db.commit()
        db.refresh(new_book)
        return new_book


@router.get("/{id}", response_model=BookSchema)
def get_book(id: int):
    with SessionLocal() as db:
        book = db.query(Book).filter(Book.id == id).first()
        if not book:
            raise HTTPException(status_code=404, detail="Book not found")
        return book


@router.put("/{id}", response_model=BookSchema)
def update_book(id: int, book: BookSchema):
    with SessionLocal() as db:
        db_book = db.query(Book).filter(Book.id == id).first()
        if not db_book:
            raise HTTPException(status_code=404, detail="Book not found")
        db_book.title = book.title
        db_book.genre = book.genre
        db_book.year = book.year
        db.commit()
        db.refresh(db_book)
        return db_book


@router.delete("/{id}", status_code=204)
def delete_book(id: int):
    with SessionLocal() as db:
        book = db.query(Book).filter(Book.id == id).first()
        if not book:
            raise HTTPException(status_code=404, detail="Book not found")
        db.delete(book)
        db.commit()
        return None
