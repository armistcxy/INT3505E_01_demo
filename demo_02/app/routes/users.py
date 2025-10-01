from flask import Blueprint, jsonify, request
from ..models import User
from ..extensions import db

users_bp = Blueprint("users", __name__)

@users_bp.route("/", methods=["GET"])
def get_users():
    users = User.query.all()
    return jsonify([{"id": u.id, "name": u.name, "email": u.email} for u in users])

@users_bp.route("/<int:user_id>", methods=["GET"])
def get_user(user_id):
    user = User.query.get_or_404(user_id)
    return jsonify({"id": user.id, "name": user.name, "email": user.email})

@users_bp.route("/", methods=["POST"])
def add_user():
    data = request.json
    user = User(name=data["name"], email=data["email"])
    db.session.add(user)
    db.session.commit()
    return jsonify({"message": "User added", "id": user.id}), 201

@users_bp.route("/<int:user_id>", methods=["PUT"])
def update_user(user_id):
    data = request.json
    user = User.query.get_or_404(user_id)
    
    if "name" in data and data["name"]:
        user.name = data["name"]
    
    if "email" in data and data["email"]:
        user.email = data["email"]
        
    db.session.commit()
    return jsonify({"message": "User updated"})

@users_bp.route("/<int:user_id>", methods=["DELETE"])
def delete_user(user_id):
    user = User.query.get_or_404(user_id)
    db.session.delete(user)
    db.session.commit()
    return jsonify({"message": "User deleted"})