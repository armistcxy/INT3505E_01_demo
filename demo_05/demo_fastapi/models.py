# models.py
from sqlalchemy import Column, Integer, String, ForeignKey, Date
from sqlalchemy.orm import relationship
from database import Base  # <-- import Base từ database.py

class Book(Base):
    __tablename__ = "books"
    id = Column(Integer, primary_key=True, index=True)
    title = Column(String)
    genre = Column(String)
    year = Column(Integer)
    authors = relationship("Author", secondary="book_authors", back_populates="books")


class Author(Base):
    __tablename__ = "authors"
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String)
    books = relationship("Book", secondary="book_authors", back_populates="authors")


class BookAuthor(Base):
    __tablename__ = "book_authors"
    book_id = Column(Integer, ForeignKey("books.id"), primary_key=True)
    author_id = Column(Integer, ForeignKey("authors.id"), primary_key=True)


class Member(Base):
    __tablename__ = "members"
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String)
    email = Column(String)
    loans = relationship("Loan", back_populates="member")


class Loan(Base):
    __tablename__ = "loans"
    id = Column(Integer, primary_key=True, index=True)
    member_id = Column(Integer, ForeignKey("members.id"))
    book_id = Column(Integer, ForeignKey("books.id"))
    loan_date = Column(Date)
    member = relationship("Member", back_populates="loans")
