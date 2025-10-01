from flask import Blueprint, jsonify, request
from ..models import Book
from ..extensions import db

books_bp = Blueprint("books", __name__)

@books_bp.route("/", methods=["GET"])
def get_books():
    books = Book.query.all()
    return jsonify([{"id": b.id, "title": b.title, "author": b.author} for b in books])

@books_bp.route("/<int:book_id>", methods=["GET"])
def get_book(book_id):
    book = Book.query.get_or_404(book_id)
    return jsonify({"id": book.id, "title": book.title, 
                    "author": book.author,
                    "publication_year": book.publication_year,
                    "genre": book.genre
                    })

@books_bp.route("/", methods=["POST"])
def add_book():
    data = request.json
    book = Book(title=data["title"], 
                author=data.get("author"),
                publication_year=data.get("publication_year"), 
                genre=data.get("genre")
                )
    db.session.add(book)
    db.session.commit()
    return jsonify({"message": "Book added", "id": book.id}), 201

@books_bp.route("/<int:book_id>", methods=["PUT"])
def update_book(book_id):
    data = request.json
    book = Book.query.get_or_404(book_id)
   
    if "title" in data and data["title"]:  
        book.title = data["title"]
    
    if "author" in data and data["author"]:  
        book.author = data["author"]
    
    if "publication_year" in data and data["publication_year"]:  
        book.publication_year = int(data["publication_year"])
    
    if "genre" in data and data["genre"]:  
        book.genre = data["genre"]
   
    db.session.commit()
    return jsonify({"message": "Book updated"})

@books_bp.route("/<int:book_id>", methods=["DELETE"])
def delete_book(book_id):
    book = Book.query.get_or_404(book_id)
    db.session.dselete(book)
    db.session.commit()
    return jsonify({"message": "Book deleted"})
    