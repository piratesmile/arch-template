package user

import (
	"arch-template/internal/app/entity"
	"arch-template/pkg/auth"
	"arch-template/pkg/errdefs"
	"context"
	"errors"
)

type Service struct {
	repo       *Repository
	tokenMaker auth.TokenManager
}

func NewService(repo *Repository, tokenMaker auth.TokenManager) *Service {
	return &Service{repo: repo, tokenMaker: tokenMaker}
}

type LoginReq struct {
	Username string
	Password string
}

type LoginResp struct {
	Token    string
	ExpireAt int64
}

func (s *Service) Login(ctx context.Context, request LoginReq) (resp LoginResp, err error) {
	user, err := s.repo.FetchByUserName(ctx, request.Username)
	if err != nil {
		return
	}
	if err = auth.ComparePassword(user.Password, request.Password); err != nil {
		return resp, errdefs.ErrIncorrectPassword
	}
	token, expireAt, err := s.tokenMaker.Generate(user.ID)
	if err != nil {
		return resp, err
	}
	return LoginResp{
		Token:    token,
		ExpireAt: expireAt,
	}, nil
}

type RegisterReq struct {
	Username        string
	Password        string
	PasswordConfirm string
}

type RegisterResp struct {
}

func (s *Service) Register(ctx context.Context, req RegisterReq) (resp RegisterResp, err error) {
	if exists, err := s.userExists(ctx, req.Username); err != nil {
		return resp, err
	} else if exists {
		return resp, errdefs.ErrUserAlreadyExists
	}
	hashPassword, err := auth.Encrypt(req.Password)
	if err != nil {
		return resp, err
	}
	user := entity.User{
		UserName: req.Username,
		Password: hashPassword,
	}
	err = s.repo.Create(ctx, &user)
	return
}

func (s *Service) userExists(ctx context.Context, username string) (bool, error) {
	user, err := s.repo.FetchByUserName(ctx, username)
	if err != nil && !errors.Is(err, errdefs.ErrResourceNotFound) {
		return false, err
	}
	return user != nil, nil
}
