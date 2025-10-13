from fastapi import FastAPI
from database import Base, engine  # <-- import đúng chỗ
from routers import books, authors, members

# Tạo bảng trong database nếu chưa có
Base.metadata.create_all(bind=engine)

app = FastAPI(title="Library API", version="1.0")

# Gắn router cho từng resource
app.include_router(books.router)
app.include_router(authors.router)
app.include_router(members.router)
