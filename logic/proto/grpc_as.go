package proto

import (
	"NewPhoto/db"
	"context"
)

type Authentication struct {
	DBInstanse db.IDB
}

func (a *Authentication) LoginUser(ctx context.Context, r *UserLoginRequest) (*UserLoginResponse, error) {
	userid, err := a.DBInstanse.LoginUser(r.GetLogin(), r.GetPassword())
	var ok bool = true
	if err != nil {
		ok = false
	}
	return &UserLoginResponse{Userid: userid, Ok: ok}, nil
}

func (a *Authentication) RegisterUser(ctx context.Context, r *UserRegisterRequest) (*UserRegisterResponse, error) {
	ok := a.DBInstanse.RegisterUser(r.GetLogin(), r.GetPassword(), r.GetFirstname(), r.GetSecondname())
	return &UserRegisterResponse{Ok: ok}, nil
}

func (a *Authentication) mustEmbedUnimplementedAuthenticationServer() {}

func NewAuthentication() *Authentication {
	r := new(Authentication)
	r.DBInstanse = db.New()
	return r
}
