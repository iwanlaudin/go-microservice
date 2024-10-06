package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/dto/response"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/models"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/repository"
	"github.com/jmoiron/sqlx"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
	DB             *sqlx.DB
}

func NewAuthService(authRepository repository.AuthRepository, db *sqlx.DB) AuthService {
	return &AuthServiceImpl{
		AuthRepository: authRepository,
		DB:             db,
	}
}

func (service *AuthServiceImpl) Create(ctx context.Context, request request.CreateUserRequest) (*response.UserResponse, error) {
	var userResponse response.UserResponse

	userId, err := uuid.NewV7()
	helpers.PanicIfError(err)

	salt := helpers.HashString("N")
	password := helpers.GeneratePasswordHash(salt, request.Password)

	now := time.Now().UTC()
	idStr := userId.String()

	user := models.User{
		ID:        userId,
		FirstName: request.FirstName,
		LastName:  &request.LastName,
		Username:  request.Email,
		Email:     request.Email,
		Salt:      salt,
		Password:  password,
		CreatedAt: &now,
		CreatedBy: &idStr,
	}

	err = service.AuthRepository.AddUser(ctx, service.DB, &user)
	helpers.PanicIfError(err)

	userResponse = response.UserResponse{
		UserId: userId,
	}

	return &userResponse, nil
}

func (service *AuthServiceImpl) SignIn(ctx context.Context, request *request.SignInRequest) (*response.UserTokenResponse, error) {
	var userTokenResponse response.UserTokenResponse

	user, err := service.AuthRepository.FindByUsername(ctx, service.DB, request.Username)
	if err != nil {
		return nil, helpers.CustomError("Failed to find user by username")
	}

	if user == nil {
		return nil, helpers.CustomError("Incorrect username or password")
	}

	if !helpers.VerifyPasswordHash(user.Password, request.Password, user.Salt) {
		return nil, helpers.CustomError("Incorrect username or password")
	}

	userClaim := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	token, err := api.GenerateToken(userClaim)
	if err != nil {
		return nil, helpers.CustomError("Failed to generate access token")
	}

	now := time.Now().UTC()

	refreshToken := helpers.GenerateRefreshToken()
	expiryAt := now.Add(14 * 24 * time.Hour)

	userTokenResponse = response.UserTokenResponse{
		UserId:       user.ID,
		Expiry:       float64(expiryAt.Unix()),
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	idStr := user.ID.String()
	userTokenId, err := uuid.NewV7()
	helpers.PanicIfError(err)
	userToken := models.UserToken{
		ID:           userTokenId,
		UserId:       user.ID,
		RefreshToken: refreshToken,
		ExpiryAt:     expiryAt,
		CreatedAt:    &now,
		CreatedBy:    &idStr,
	}

	err = service.AuthRepository.AddUserToken(ctx, service.DB, &userToken)
	if err != nil {
		return nil, helpers.CustomError("Failed to add user token")
	}

	return &userTokenResponse, nil
}

func (service *AuthServiceImpl) RefreshToken(ctx context.Context, request *request.RefreshTokenRequest) (*response.UserTokenResponse, error) {
	var userTokenResponse response.UserTokenResponse
	userContext := api.UserFromContext(ctx)

	err := helpers.WithTransaction(ctx, service.DB, func(tx *sqlx.Tx) error {
		userToken, err := service.AuthRepository.FindUserTokenByRefreshToken(ctx, service.DB, request.RefreshToken)
		if err != nil {
			return helpers.CustomError("Failed to find user token")
		}

		if userToken == nil || userToken.IsUsed || userToken.ExpiryAt.Before(time.Now().UTC()) {
			return helpers.CustomError("Invalid or expired refresh token")
		}

		now := time.Now().UTC()
		userToken.UpdatedAt = &now
		userToken.UpdatedBy = &userContext.ID

		err = service.AuthRepository.UpdateUserToken(ctx, service.DB, userToken)
		if err != nil {
			return helpers.CustomError("Failed to update user token")
		}

		refreshToken := helpers.GenerateRefreshToken()
		expiryAt := time.Now().UTC().Add(14 * 24 * time.Hour)
		userTokenId, err := uuid.NewV7()
		if err != nil {
			return helpers.CustomError("Failed to generate token ID")
		}

		newUserToken := models.UserToken{
			ID:           userTokenId,
			UserId:       userToken.UserId,
			ExpiryAt:     expiryAt,
			RefreshToken: refreshToken,
			CreatedBy:    &userContext.ID,
			CreatedAt:    &now,
		}

		err = service.AuthRepository.AddUserToken(ctx, service.DB, &newUserToken)
		if err != nil {
			return helpers.CustomError("Failed to add new user token")
		}

		user, err := service.AuthRepository.FindById(ctx, service.DB, userToken.UserId)
		if err != nil {
			return helpers.CustomError("Failed to find user")
		}
		if user == nil {
			return helpers.CustomError("User not found")
		}

		userClaim := map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		}

		token, err := api.GenerateToken(userClaim)
		if err != nil {
			return helpers.CustomError("Failed to generate access token")
		}

		userTokenResponse = response.UserTokenResponse{
			UserId:       user.ID,
			Expiry:       float64(expiryAt.Unix()),
			AccessToken:  token,
			RefreshToken: refreshToken,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &userTokenResponse, nil
}

func (service *AuthServiceImpl) FindUserById(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
	var userResponse response.UserResponse

	user, err := service.AuthRepository.FindById(ctx, service.DB, id)
	if err != nil {
		return nil, helpers.CustomError("Failed to find user")
	}
	if user == nil {
		return nil, helpers.CustomError("User not found")
	}

	userResponse = response.UserResponse{
		UserId:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
	}

	return &userResponse, nil
}
