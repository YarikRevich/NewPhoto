package proto

import (
	"NewPhoto/db"
	"context"
)

type Authentication struct {
	DBInstanse db.IDB
}

func (a *Authentication) LoginUser(ctx context.Context, r *UserLoginRequest) (*UserLoginResponse, error) {
	accessToken, loginToken, err := a.DBInstanse.Login(r.GetLogin(), r.GetPassword(), r.GetSourceType())
	var ok bool = true
	if err != nil {
		ok = false
	}
	return &UserLoginResponse{AccessToken: accessToken, LoginToken: loginToken, Ok: ok}, nil
}

func (a *Authentication) RegisterUser(ctx context.Context, r *UserRegisterRequest) (*UserRegisterResponse, error) {
	err := a.DBInstanse.RegisterUser(r.GetLogin(), r.GetPassword(), r.GetFirstname(), r.GetSecondname())
	if err != nil {
		return &UserRegisterResponse{Ok: false}, nil
	}
	return &UserRegisterResponse{Ok: true}, nil
}

func (a *Authentication) IsTokenCorrect(ctx context.Context, r *IsTokenCorrectRequest) (*IsTokenCorrectResponse, error) {
	ok := a.DBInstanse.IsTokenCorrect(r.GetAccessToken(), r.GetLoginToken(), r.GetSourceType())
	return &IsTokenCorrectResponse{Ok: ok}, nil
}

// func (a *Authentication) RetrieveToken(ctx context.Context, r *RetrieveTokenRequest) (*RetrieveTokenResponse, error) {
// 	accessToken, loginToken, ok := a.DBInstanse.RetrieveToken(r.GetAccessToken(), r.GetLoginToken(), r.GetSourceType())
// 	if !ok {
// 		return &RetrieveTokenResponse{AccessToken: accessToken, LoginToken: loginToken, Ok: false}, nil
// 	}
// 	return &RetrieveTokenResponse{AccessToken: accessToken, LoginToken: loginToken, Ok: true}, nil
// }

func (a *Authentication) mustEmbedUnimplementedAuthenticationServer() {}

func NewAuthentication() *Authentication {
	r := new(Authentication)
	r.DBInstanse = db.New()
	return r
}
