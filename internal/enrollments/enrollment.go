package enrollments

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Role string

const (
	RoleOwner   Role = "owner"
	RoleManager Role = "manager"
	RoleLearner Role = "learner"
)

type Enrollment struct {
	CourseID  string    `json:"course_id"`
	UserID    string    `json:"user_id"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, enrollment Enrollment) (Enrollment, error)
	Delete(ctx context.Context, courseID, userID string) error
	Get(ctx context.Context, courseID, userID string) (Enrollment, error)
	ListByCourse(ctx context.Context, courseID string) ([]Enrollment, error)
	DB() *sql.DB
	WithTx(tx *sql.Tx) Repository
}

var (
	ErrNotFound = errors.New("enrollment not found")
)
