from fastapi import APIRouter, HTTPException
from sqlalchemy.orm import Session
from typing import List, Optional
from models import Member
from pydantic import BaseModel
from database import SessionLocal

router = APIRouter(prefix="/members", tags=["Members"])

class MemberSchema(BaseModel):
    id: Optional[int]
    name: str
    email: str

    class Config:
        orm_mode = True


class LoanSchema(BaseModel):
    id: Optional[int]
    member_id: int
    book_id: int
    class Config:
        orm_mode = True


@router.get("/", response_model=List[MemberSchema])
def list_members(
    name: Optional[str] = None,
    email: Optional[str] = None,
    limit: int = 10,
    offset: int = 0,
):
    with SessionLocal() as db:
        query = db.query(Member)
        if name:
            query = query.filter(Member.name.ilike(f"%{name}%"))
        if email:
            query = query.filter(Member.email.ilike(f"%{email}%"))
        return query.order_by(Member.id).limit(limit).offset(offset).all()


@router.post("/", response_model=MemberSchema, status_code=201)
def create_member(member: MemberSchema):
    with SessionLocal() as db:
        new_member = Member(name=member.name, email=member.email)
        db.add(new_member)
        db.commit()
        db.refresh(new_member)
        return new_member


@router.get("/{id}", response_model=MemberSchema)
def get_member(id: int):
    with SessionLocal() as db:
        member = db.query(Member).filter(Member.id == id).first()
        if not member:
            raise HTTPException(status_code=404, detail="Member not found")
        return member


@router.get("/{id}/loans", response_model=List[LoanSchema])
def get_loans_by_member(id: int):
    with SessionLocal() as db:
        member = db.query(Member).filter(Member.id == id).first()
        if not member:
            raise HTTPException(status_code=404, detail="Member not found")
        return member.loans
