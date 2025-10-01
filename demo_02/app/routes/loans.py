from flask import Blueprint, jsonify, request
from ..models import Loan, Book, User
from ..extensions import db
from datetime import datetime

loans_bp = Blueprint("loans", __name__)

@loans_bp.route("/", methods=["GET"])
def get_loans():
    loans = Loan.query.all()
    return jsonify([{
        "loan_id": l.loan_id, 
        "book_id": l.book_id,
        "book_title": l.book.title,
        "user_id": l.user_id,
        "user_name": l.user.name,
        "checkout_date": l.checkout_date.isoformat(),
        "return_date": l.return_date.isoformat() if l.return_date else None
    } for l in loans])

@loans_bp.route("/<int:loan_id>", methods=["GET"])
def get_loan(loan_id):
    loan = Loan.query.get_or_404(loan_id)
    return jsonify({
        "loan_id": loan.loan_id,
        "book_id": loan.book_id,
        "book_title": loan.book.title,
        "user_id": loan.user_id,
        "user_name": loan.user.name,
        "checkout_date": loan.checkout_date.isoformat(),
        "return_date": loan.return_date.isoformat() if loan.return_date else None
    })

@loans_bp.route("/", methods=["POST"])
def add_loan():
    data = request.json
    book_id = data.get("book_id")
    user_id = data.get("user_id")

    if not book_id or not user_id:
        return jsonify({"message": "Book ID and User ID are required"}), 400

    book = Book.query.get(book_id)
    if not book:
        return jsonify({"message": "Book not found"}), 404
    
    if book.count <= 0:
        return jsonify({"message": "Book is out of stock"}), 400

    user = User.query.get(user_id)
    if not user:
        return jsonify({"message": "User not found"}), 404

    book.count -= 1

    loan = Loan(
        book_id=book_id,
        user_id=user_id,
        checkout_date=datetime.utcnow().date()
    )
    db.session.add(loan)
    db.session.commit()
    return jsonify({"message": "Loan added", "loan_id": loan.loan_id}), 201

@loans_bp.route("/<int:loan_id>", methods=["PUT"])
def update_loan(loan_id):
    loan = Loan.query.get_or_404(loan_id)
    
    if not loan.return_date:
        loan.return_date = datetime.utcnow().date()
        
        book = Book.query.get(loan.book_id)
        if book:
            book.count += 1
            
        db.session.commit()
        return jsonify({"message": "Loan updated with return date"})
    else:
        return jsonify({"message": "Book has already been returned"}), 400


@loans_bp.route("/<int:loan_id>", methods=["DELETE"])
def delete_loan(loan_id):
    loan = Loan.query.get_or_404(loan_id)

    if not loan.return_date:
        book = Book.query.get(loan.book_id)
        if book:
            book.count += 1

    db.session.delete(loan)
    db.session.commit()
    return jsonify({"message": "Loan deleted"})