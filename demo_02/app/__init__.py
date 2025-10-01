from flask import Flask
from .extensions import db
from .routes.books import books_bp
from .routes.users import users_bp
from .routes.loans import loans_bp

def create_app():
    app = Flask(__name__)
    app.config.from_object("config.Config")
    
    db.init_app(app)
    
    app.register_blueprint(books_bp, url_prefix="/books")
    app.register_blueprint(users_bp, url_prefix="/users")
    app.register_blueprint(loans_bp, url_prefix="/loans")
    
    return app