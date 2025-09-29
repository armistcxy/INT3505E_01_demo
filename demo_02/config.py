import os

class Config:
    # Flask configuration
    DEBUG = True
    SECRET_KEY = os.environ.get("SECRET_KEY", "mysecretkey")

    # SQLAlchemy Configuration
    DB_USER = os.environ.get("DB_USER", "myuser")
    DB_PASSWORD = os.environ.get("DB_PASSWORD", "mypassword")
    DB_HOST = os.environ.get("DB_HOST", "localhost")
    DB_PORT = os.environ.get("DB_PORT", "5432")
    DB_NAME = os.environ.get("DB_NAME", "mydatabase")

    SQLALCHEMY_DATABASE_URI = f"postgresql+psycopg2://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"
    SQLALCHEMY_TRACK_MODIFICATIONS = False