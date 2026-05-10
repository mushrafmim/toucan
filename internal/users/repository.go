package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"toucan/internal/database"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u User) (User, error) {
	if u.ID == "" {
		u.ID = database.NewID()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	u.UpdatedAt = u.CreatedAt

	rolesJSON, err := json.Marshal(u.Roles)
	if err != nil {
		return User{}, fmt.Errorf("marshal roles: %w", err)
	}

	query := `
		INSERT INTO users (id, external_subject, email, display_name, roles, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = r.db.ExecContext(ctx, query, u.ID, u.ExternalSubject, u.Email, u.DisplayName, rolesJSON, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return User{}, fmt.Errorf("postgres create user: %w", err)
	}
	return u, nil
}

func (r *Repository) Get(ctx context.Context, id string) (User, error) {
	query := `
		SELECT id, external_subject, email, display_name, roles, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	var rolesJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.ExternalSubject, &u.Email, &u.DisplayName, &rolesJSON, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("postgres get user: %w", err)
	}

	if err := json.Unmarshal(rolesJSON, &u.Roles); err != nil {
		return User{}, fmt.Errorf("unmarshal roles: %w", err)
	}

	return u, nil
}

func (r *Repository) GetByExternalSubject(ctx context.Context, subject string) (User, error) {
	query := `
		SELECT id, external_subject, email, display_name, roles, created_at, updated_at
		FROM users
		WHERE external_subject = $1
	`
	var u User
	var rolesJSON []byte
	err := r.db.QueryRowContext(ctx, query, subject).Scan(&u.ID, &u.ExternalSubject, &u.Email, &u.DisplayName, &rolesJSON, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("postgres get user by subject: %w", err)
	}

	if err := json.Unmarshal(rolesJSON, &u.Roles); err != nil {
		return User{}, fmt.Errorf("unmarshal roles: %w", err)
	}

	return u, nil
}

func (r *Repository) List(ctx context.Context, filter ListFilter) (ListResult, error) {
	page := max(filter.Page, 1)
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total); err != nil {
		return ListResult{Page: page, PageSize: pageSize}, nil
	}

	query := `
		SELECT id, external_subject, email, display_name, roles, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return ListResult{Page: page, PageSize: pageSize, Total: total}, fmt.Errorf("postgres list users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		var rolesJSON []byte
		if err := rows.Scan(&u.ID, &u.ExternalSubject, &u.Email, &u.DisplayName, &rolesJSON, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return ListResult{Page: page, PageSize: pageSize, Total: total}, fmt.Errorf("postgres scan user: %w", err)
		}

		if err := json.Unmarshal(rolesJSON, &u.Roles); err != nil {
			return ListResult{Page: page, PageSize: pageSize, Total: total}, fmt.Errorf("unmarshal roles for user %s: %w", u.ID, err)
		}

		users = append(users, u)
	}

	return ListResult{
		Items:    users,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (r *Repository) Update(ctx context.Context, u User) (User, error) {
	u.UpdatedAt = time.Now().UTC()
	rolesJSON, err := json.Marshal(u.Roles)
	if err != nil {
		return User{}, fmt.Errorf("marshal roles: %w", err)
	}

	query := `
		UPDATE users
		SET email = $2, display_name = $3, roles = $4, updated_at = $5
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.DisplayName, rolesJSON, u.UpdatedAt)
	if err != nil {
		return User{}, fmt.Errorf("postgres update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return User{}, err
	}
	if rows == 0 {
		return User{}, ErrNotFound
	}

	return u, nil
}
