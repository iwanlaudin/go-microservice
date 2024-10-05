package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindById(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*models.User, error)
	FindByUsername(ctx context.Context, db *sqlx.DB, username string) (*models.User, error)
	FindByEmail(ctx context.Context, db *sqlx.DB, email string) (*models.User, error)
	AddUser(ctx context.Context, db *sqlx.DB, user *models.User) error
	AddUserToken(ctx context.Context, db *sqlx.DB, userToken *models.UserToken) error
	FindUserTokenByRefreshToken(ctx context.Context, db *sqlx.DB, refreshToken string) (*models.UserToken, error)
}
