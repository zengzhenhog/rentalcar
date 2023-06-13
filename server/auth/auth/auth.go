package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	OpenIdResolver OpenIdResolver
	Mongo          *dao.Mongo
	Logger         zap.Logger
	TokenGenerator TokenGenerator
	TokenExpire    time.Duration
}

type OpenIdResolver interface {
	Resolve(code string) (string, error)
}

type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error) // accountID string 跟具体业务不太相关，为了更通用写成这样
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// s.Logger.Info("received code", zap.String("code", req.Code))
	openID, err := s.OpenIdResolver.Resolve(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openid: %v", err)
	}

	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		s.Logger.Error("account resolve account id", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	// tkn, err := s.TokenGenerator.GenerateToken(accountID, s.TokenExpire)
	tkn, err := s.TokenGenerator.GenerateToken(accountID.String(), s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: tkn,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}
