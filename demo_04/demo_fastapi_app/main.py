from fastapi import FastAPI, HTTPException, Depends, status, Response
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from sqlmodel import SQLModel, Field, Session, create_engine, select
from typing import Optional, List
from uuid import uuid4
from jose import JWTError, jwt
from passlib.context import CryptContext
from datetime import datetime, timedelta
from fastapi import Security
from fastapi.security import OAuth2PasswordBearer, SecurityScopes



class User(SQLModel, table=True):
    id: Optional[str] = Field(default_factory=lambda: str(uuid4()), primary_key=True)
    username: str
    hashed_password: str

class UserCreate(SQLModel):
    username: str
    password: str

class Token(SQLModel):
    access_token: str
    token_type: str = "bearer"

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

SECRET_KEY = "secret_key_demo"
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="users/login")


SQLITE_FILE_NAME = "movies.db"
DATABASE_URL = f"sqlite:///{SQLITE_FILE_NAME}"
engine = create_engine(DATABASE_URL, echo=False)

def create_db_and_tables():
    SQLModel.metadata.create_all(engine)

def get_user_by_username(username: str) -> Optional[User]:
    with Session(engine) as session:
        stmt = select(User).where(User.username == username)
        return session.exec(stmt).first()


def verify_password(plain, hashed):
    return pwd_context.verify(plain, hashed)


def hash_password(password):
    return pwd_context.hash(password)


def create_access_token(data: dict, expires_delta: Optional[timedelta] = None):
    to_encode = data.copy()
    expire = datetime.utcnow() + (expires_delta or timedelta(minutes=15))
    to_encode.update({"exp": expire})
    return jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)

def get_current_user(token: str = Depends(oauth2_scheme)):
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Invalid authentication credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        if username is None:
            raise credentials_exception
    except JWTError:
        raise credentials_exception

    user = get_user_by_username(username)
    if user is None:
        raise credentials_exception
    return user

# -----------------------------
# FastAPI app
# -----------------------------
app = FastAPI(title="Movie API", version="1.0.0")

@app.post("/users/register", response_model=dict)
def register_user(user_in: UserCreate):
    if get_user_by_username(user_in.username):
        raise HTTPException(status_code=400, detail="Username already registered")

    user = User(username=user_in.username, hashed_password=hash_password(user_in.password))
    with Session(engine) as session:
        session.add(user)
        session.commit()
    return {"msg": "User registered successfully"}

@app.post("/users/login", response_model=Token)
def login_user(form_data: OAuth2PasswordRequestForm = Depends()):
    user = get_user_by_username(form_data.username)
    if not user or not verify_password(form_data.password, user.hashed_password):
        raise HTTPException(status_code=400, detail="Incorrect username or password")

    access_token = create_access_token(data={"sub": user.username}, expires_delta=timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES))
    return Token(access_token=access_token)

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


@app.get("/movies", response_model=List[Movie])
def get_all_movies():
    """Get all movies."""
    with Session(engine) as session:
        movies = session.exec(select(Movie)).all()
        return movies

@app.post(
    "/movies",
    response_model=Movie,
    status_code=status.HTTP_201_CREATED,
)
def create_movie(movie_in: MovieCreate, 
                 response: Response,
                 current_user: User = Security(get_current_user),):
    movie = Movie.from_orm(movie_in)
    with Session(engine) as session:
        session.add(movie)
        session.commit()
        session.refresh(movie)

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