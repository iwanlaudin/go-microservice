package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepositoryImpl struct {
}

func NewAuthRepository() AuthRepository {
	return &AuthRepositoryImpl{}
}

func (repository *AuthRepositoryImpl) FindById(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*models.User, error) {
	user := models.User{}
	err := db.GetContext(ctx, &user, `SELECT * FROM "user" WHERE is_deleted=false AND id=$1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		helpers.PanicIfError(err)
	}

	return &user, nil
}

func (repositpry *AuthRepositoryImpl) FindByUsername(ctx context.Context, db *sqlx.DB, username string) (*models.User, error) {
	user := models.User{}
	err := db.GetContext(ctx, &user, `SELECT * FROM "user" WHERE is_deleted=false AND username=$1`, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		helpers.PanicIfError(err)
	}

	return &user, nil
}

func (repositpry *AuthRepositoryImpl) FindByEmail(ctx context.Context, db *sqlx.DB, email string) (*models.User, error) {
	user := models.User{}
	err := db.GetContext(ctx, &user, `SELECT * FROM "user" WHERE is_deleted=false AND email=$1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		helpers.PanicIfError(err)
	}

	return &user, nil
}

func (repository *AuthRepositoryImpl) AddUser(ctx context.Context, db *sqlx.DB, user *models.User) error {
	query := `
		INSERT INTO "user" (
			id, first_name, last_name, username, email, salt, password, created_at, created_by
		) VALUES (
		 	:id, :first_name, :last_name, :username, :email, :salt, :password, :created_at, :created_by
		)
	`
	_, err := db.NamedExecContext(ctx, query, user)
	helpers.PanicIfError(err)

	return nil
}

func (repository *AuthRepositoryImpl) AddUserToken(ctx context.Context, db *sqlx.DB, user *models.UserToken) error {
	query := `
		INSERT INTO "userToken" (
			id, user_id, refresh_token, expiry_at, created_at, created_by
		) VALUES (
		 	:id, :user_id, :refresh_token, :expiry_at, :created_at, :created_by
		)
	`
	_, err := db.NamedExecContext(ctx, query, user)
	helpers.PanicIfError(err)

	return nil
}

func (repository *AuthRepositoryImpl) UpdateUserToken(ctx context.Context, db *sqlx.DB, userToken *models.UserToken) error {
	query := `
		UPDATE "userToken"
		SET
			used_at = :used_at, 
            is_used = :is_used, 
            updated_at = :updated_at,
			updated_by = :updated_by
		WHERE
			is_deleted = :is_deleted
		AND
			id = :id
	`
	now := time.Now().UTC()
	_, err := db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         userToken.ID,
		"used_at":    now,
		"is_used":    true,
		"updated_at": now,
		"updated_by": userToken.UpdatedBy,
		"is_deleted": false,
	})
	helpers.PanicIfError(err)

	return nil
}

func (repository *AuthRepositoryImpl) FindUserTokenByRefreshToken(ctx context.Context, db *sqlx.DB, refreshToken string) (*models.UserToken, error) {
	var userToken models.UserToken

	query := `
		SELECT 
			ut.id AS id,
			ut.user_id AS user_id,
			ut.refresh_token AS refresh_token,
			ut.is_used AS is_used,
			ut.expiry_at AS expiry_at,
			ut.is_deleted AS is_deleted
		FROM 
			"userToken" ut
		WHERE 
			ut.refresh_token = $1
		AND
			ut.is_deleted = false
	`
	err := db.GetContext(ctx, &userToken, query, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		helpers.PanicIfError(err)
	}

	return &userToken, nil
}
