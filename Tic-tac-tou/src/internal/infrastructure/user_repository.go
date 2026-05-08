package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"tic-tac/internal/infrastructure/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, login, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	query := `INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3) RETURNING id, login, created_at`
	var user model.User
	err = r.db.QueryRow(ctx, query, id, login, string(hash)).Scan(&user.ID, &user.Login, &user.CreatedAt)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_login_key\" (SQLSTATE 23505)" {
			return nil, fmt.Errorf("user with this login already exists")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	row := r.db.QueryRow(ctx, query, login)
	var user model.User
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ValidatePassword(user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
