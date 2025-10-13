from typing import Optional, List
from uuid import uuid4

from fastapi import FastAPI, HTTPException, status, Response
from sqlmodel import SQLModel, Field, create_engine, Session, select

# -----------------------------
# Models
# -----------------------------
class Movie(SQLModel, table=True):
    """Database model + Pydantic model for Movie.

    - id: UUID string (primary key)
    - title: required
    - description: optional
    - year: required
    - genre: optional
    """
    id: Optional[str] = Field(default_factory=lambda: str(uuid4()), primary_key=True)
    title: str
    description: Optional[str] = None
    year: int
    genre: Optional[str] = None

# Schema used for create/update requests (exclude id)
class MovieCreate(SQLModel):
    title: str
    description: Optional[str] = None
    year: int
    genre: Optional[str] = None

# -----------------------------
# Database setup
# -----------------------------
SQLITE_FILE_NAME = "movies.db"
DATABASE_URL = f"sqlite:///{SQLITE_FILE_NAME}"
engine = create_engine(DATABASE_URL, echo=False)

# Create tables at startup
def create_db_and_tables():
    SQLModel.metadata.create_all(engine)

# -----------------------------
# FastAPI app
# -----------------------------
app = FastAPI(title="Movie API", version="1.0.0")

@app.on_event("startup")
def on_startup():
    # ensure database + tables exist
    create_db_and_tables()

    # optional: seed with a few example movies if DB empty (handy for demo)
    with Session(engine) as session:
        stmt = select(Movie)
        count = session.exec(stmt).all()
        if not count:
            demo = [
                Movie(title="Inception", description="A thief who steals corporate secrets through dream-sharing technology.", year=2010, genre="Sci-Fi"),
                Movie(title="The Matrix", description="A hacker discovers reality is a simulation.", year=1999, genre="Sci-Fi"),
                Movie(title="Her", description="A man falls in love with an intelligent operating system.", year=2013, genre="Drama"),
            ]
            session.add_all(demo)
            session.commit()

# -----------------------------
# Endpoints (match your OpenAPI spec)
# -----------------------------

@app.get("/movies", response_model=List[Movie])
def get_all_movies():
    """Get all movies."""
    with Session(engine) as session:
        movies = session.exec(select(Movie)).all()
        return movies

@app.post("/movies", response_model=Movie, status_code=status.HTTP_201_CREATED)
def create_movie(movie_in: MovieCreate, response: Response):
    """Create a new movie. Returns the created movie (201)."""
    movie = Movie.from_orm(movie_in)  # will generate id
    with Session(engine) as session:
        session.add(movie)
        session.commit()
        session.refresh(movie)

    # set Location header pointing to the new resource
    response.headers["Location"] = f"/movies/{movie.id}"
    return movie

@app.get("/movies/{id}", response_model=Movie)
def get_movie(id: str):
    """Get movie by ID. 404 if not found."""
    with Session(engine) as session:
        movie = session.get(Movie, id)
        if not movie:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Movie not found")
        return movie

@app.put("/movies/{id}", response_model=Movie)
def update_movie(id: str, movie_in: MovieCreate):
    """Update movie by ID. 404 if not found."""
    with Session(engine) as session:
        movie = session.get(Movie, id)
        if not movie:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Movie not found")
        movie.title = movie_in.title
        movie.description = movie_in.description
        movie.year = movie_in.year
        movie.genre = movie_in.genre
        session.add(movie)
        session.commit()
        session.refresh(movie)
        return movie

@app.delete("/movies/{id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_movie(id: str):
    """Delete movie by ID. 404 if not found. Returns 204 no content on success."""
    with Session(engine) as session:
        movie = session.get(Movie, id)
        if not movie:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Movie not found")
        session.delete(movie)
        session.commit()
        return Response(status_code=status.HTTP_204_NO_CONTENT)

# -----------------------------
# Run with: uvicorn main:app --reload
# -----------------------------
