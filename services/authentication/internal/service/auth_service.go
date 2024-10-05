package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/dto/response"
)

type AuthService interface {
	Create(ctx context.Context, request request.CreateUserRequest) (*response.UserResponse, error)
	SignIn(ctx context.Context, request *request.SignInRequest) (*response.UserTokenResponse, error)
	RefreshToken(ctx context.Context, request *request.RefreshTokenRequest) (*response.UserTokenResponse, error)
	FindUserById(ctx context.Context, id uuid.UUID) (*response.UserResponse, error)
}
