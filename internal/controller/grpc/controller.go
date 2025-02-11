package grpc

import (
	"context"
	"github.com/watchlist-kata/auth/internal/service"

	"github.com/watchlist-kata/auth/internal/controller/grpc/proto"
)

type AuthServiceGRPC struct {
	proto.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthServiceGRPC(authService service.AuthService) *AuthServiceGRPC {
	return &AuthServiceGRPC{authService: authService}
}

func (a *AuthServiceGRPC) Login(ctx context.Context, r *proto.LoginRequest) (*proto.LoginResponse, error) {
	//TODO добавить обращение к сервису user с проверкой входных данных если true движемся дальше
	//попутно получаем userId
	userId := "4235"
	JWT, refreshToken, err := a.authService.GenerateTokens(ctx, userId, r.Email)
	if err != nil {
		return nil, err
	}

	return &proto.LoginResponse{Success: true, AccessToken: JWT, RefreshToken: refreshToken}, nil

}

func (a *AuthServiceGRPC) RefreshToken(ctx context.Context, r *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	jwt, err := a.authService.RefreshToken(ctx, r.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &proto.RefreshTokenResponse{Success: true, AccessToken: jwt}, nil
}

func (a *AuthServiceGRPC) ValidateJWT(ctx context.Context, r *proto.ValidateJWTRequest) (*proto.ValidateJWTResponse, error) {
	err := a.authService.ValidateJWT(r.AccessToken)
	if err != nil {
		return nil, err
	}
	return &proto.ValidateJWTResponse{Success: true}, nil
}
