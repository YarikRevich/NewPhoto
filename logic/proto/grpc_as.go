package proto

import (
	"github.com/YarikRevich/NewPhoto/db"
	"context"
)

type Authentication struct {
	DBInstanse db.IDB
}

func (a *Authentication) LoginUser(ctx context.Context, r *UserLoginRequest) (*UserLoginResponse, error) {

	accessToken, loginToken, err := a.DBInstanse.Login(r.GetLogin(), r.GetPassword(), db.SourceType(r.GetSourceType().Number()))
	var ok bool = true
	if err != nil {
		ok = false
	}
	return &UserLoginResponse{AccessToken: accessToken, LoginToken: loginToken, Ok: ok}, nil
}

func (a *Authentication) LogoutUser(ctx context.Context, r *UserLogoutRequest) (*UserLogoutResponse, error) {
	a.DBInstanse.Logout(a.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), db.SourceType(r.GetSourceType().Number()))
	return &UserLogoutResponse{Ok: true}, nil
}

func (a *Authentication) RegisterUser(ctx context.Context, r *UserRegisterRequest) (*UserRegisterResponse, error) {
	err := a.DBInstanse.RegisterUser(r.GetLogin(), r.GetPassword(), r.GetFirstname(), r.GetSecondname())
	if err != nil {
		return &UserRegisterResponse{Ok: false}, nil
	}
	return &UserRegisterResponse{Ok: true}, nil
}

func (a *Authentication) IsTokenCorrect(ctx context.Context, r *IsTokenCorrectRequest) (*IsTokenCorrectResponse, error) {
	ok := a.DBInstanse.IsTokenCorrect(r.GetAccessToken(), r.GetLoginToken(), db.SourceType(r.GetSourceType().Number()))
	return &IsTokenCorrectResponse{Ok: ok}, nil
}

func (a *Authentication) mustEmbedUnimplementedAuthenticationServer() {}

func NewAuthentication() *Authentication {
	r := new(Authentication)
	r.DBInstanse = db.New()
	return r
}
