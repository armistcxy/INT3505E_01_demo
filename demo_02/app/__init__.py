from flask import Flask
from .extensions import db
from .routes.books import books_bp

def create_app():
    app = Flask(__name__)
    app.config.from_object("config.Config")
    
    db.init_app(app)
    
    app.register_blueprint(books_bp, url_prefix="/books")
    
    return app