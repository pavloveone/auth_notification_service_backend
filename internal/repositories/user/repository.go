package userrepository

import (
	"auth_notification_service/internal/auth"
	"auth_notification_service/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

var ErrUserAlreadyExists = errors.New("user already exists")

func NewUserRepository(ctx context.Context, db *pgxpool.Pool) (*UserRepository, error) {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users ( 
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(ctx, createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("an error occured while creating table: %w", err)
	}
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) AllUsers(ctx context.Context) ([]models.UserResponse, error) {
	query := `SELECT id, username, email, created_on FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []models.UserResponse{}, fmt.Errorf("failed to parse query: %w", err)
	}
	defer rows.Close()

	users := make([]models.UserResponse, 0)
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedOn)
		if err != nil {
			return []models.UserResponse{}, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return []models.UserResponse{}, fmt.Errorf("error while reading row: %w", err)
	}
	return users, nil
}

func (r *UserRepository) UserById(ctx context.Context, id int) (models.UserResponse, error) {
	query := `SELECT id, username, email, created_on FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user models.UserResponse
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedOn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserResponse{}, fmt.Errorf("failed to find user: %w", err)
		}
		return models.UserResponse{}, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *UserRepository) AddNewUser(ctx context.Context, request models.UserCreateRequest) (int, error) {
	hashPass, err := auth.HashPassword(request.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	query := `
	INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	var id int
	err = r.db.QueryRow(ctx, query, request.Username, hashPass, request.Email).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, ErrUserAlreadyExists
			}
		}
		return 0, fmt.Errorf("failed to add new user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) LogIn(ctx context.Context, request models.UserLogInRequest) (models.UserLogInResponse, error) {
	query := `SELECT id, username, password, email, created_on FROM users WHERE username = $1`
	row := r.db.QueryRow(ctx, query, request.Username)
	var user models.UserResponse
	var hashedPass string
	err := row.Scan(&user.ID, &user.Username, &hashedPass, &user.Email, &user.CreatedOn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserLogInResponse{}, fmt.Errorf("user not found: %w", err)
		}
		return models.UserLogInResponse{}, fmt.Errorf("failed to login: %w", err)
	}
	if !auth.CheckPassHash(request.Password, hashedPass) {
		return models.UserLogInResponse{}, errors.New("invalid username or password")
	}
	access, refresh, err := auth.GenerateTokens(user.ID)
	if err != nil {
		return models.UserLogInResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}
	return models.UserLogInResponse{
		User:         user,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (r *UserRepository) LogOut(ctx context.Context) (bool, error) {
	return true, nil
}
