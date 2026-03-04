package auth_service

import (
	models "auth_service/internal/models"
	repository "auth_service/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func addUserCredentials(
	ctx context.Context,
	db *sql.DB,
	user *models.UserRegisterInfo,
) (uuid.UUID, error) {
	exists, err := checkUserExistenceByLogin(ctx, db, &user.Login)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("check existence failed: %w", err)
	}

	if exists {
		return uuid.UUID{}, models.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("hashing password failed: %w", err)
	}

	userId, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("generating user id failed: %w", err)
	}

	query := `INSERT INTO auth_credentials (user_id, login, password_hash, created_at) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(query, userId, user.Login, hashedPassword, time.Now().Unix())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("inserting user credentials failed: %w", err)
	}

	return userId, nil
}

func checkUserExistenceByLogin(ctx context.Context, db *sql.DB, login *string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM auth_credentials WHERE login = $1)`
	err := db.QueryRowContext(ctx, query, login).Scan(&exists)
	if err != nil {
		return false, repository.MapError(err)
	}
	return exists, err
}
