from fastapi import APIRouter, HTTPException, Query
from sqlalchemy.orm import Session
from typing import List, Optional
from models import Author
from pydantic import BaseModel
from database import SessionLocal

router = APIRouter(prefix="/authors", tags=["Authors"])

class AuthorSchema(BaseModel):
    id: Optional[int]
    name: str

    class Config:
        orm_mode = True


class BookSchema(BaseModel):
    id: Optional[int]
    title: str
    genre: Optional[str]
    year: Optional[int]

    class Config:
        orm_mode = True


@router.get("/", response_model=List[AuthorSchema])
def list_authors(name: Optional[str] = None, limit: int = 10, offset: int = 0):
    with SessionLocal() as db:
        query = db.query(Author)
        if name:
            query = query.filter(Author.name.ilike(f"%{name}%"))
        return query.order_by(Author.id).limit(limit).offset(offset).all()


@router.post("/", response_model=AuthorSchema, status_code=201)
def create_author(author: AuthorSchema):
    with SessionLocal() as db:
        new_author = Author(name=author.name)
        db.add(new_author)
        db.commit()
        db.refresh(new_author)
        return new_author


@router.get("/{id}", response_model=AuthorSchema)
def get_author(id: int):
    with SessionLocal() as db:
        author = db.query(Author).filter(Author.id == id).first()
        if not author:
            raise HTTPException(status_code=404, detail="Author not found")
        return author


@router.get("/{id}/books", response_model=List[BookSchema])
def get_books_by_author(id: int, limit: int = 10, offset: int = 0):
    with SessionLocal() as db:
        author = db.query(Author).filter(Author.id == id).first()
        if not author:
            raise HTTPException(status_code=404, detail="Author not found")
        return author.books[offset : offset + limit]
