package users

import (
	"context"
	"time"

	"toucan/internal/shared"
)

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleInstructor Role = "instructor"
	RoleLearner    Role = "learner"
)

type User struct {
	ID              string    `json:"id"`
	ExternalSubject string    `json:"external_subject"`
	Email           string    `json:"email"`
	DisplayName     string    `json:"display_name"`
	Roles           []Role    `json:"roles"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	ExternalSubject string `json:"external_subject"`
	Email           string `json:"email"`
	DisplayName     string `json:"displayName"`
	Roles           []Role `json:"roles"`
}

type ListFilter struct {
	Page     int
	PageSize int
}

type ListResult = shared.ListResult[User]

type Service interface {
	Create(ctx context.Context, req CreateUserRequest) (User, error)
	Get(ctx context.Context, id string) (User, error)
	GetByExternalSubject(ctx context.Context, subject string) (User, error)
	List(ctx context.Context, filter ListFilter) (ListResult, error)
}
