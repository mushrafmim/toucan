package enrollments

import (
	"context"
	"database/sql"
	"time"
)

type DBOrTx interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type concreteRepository struct {
	db   DBOrTx
	conn *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &concreteRepository{db: db, conn: db}
}

func (r *concreteRepository) DB() *sql.DB {
	return r.conn
}

func (r *concreteRepository) WithTx(tx *sql.Tx) Repository {
	if tx == nil {
		return r
	}
	return &concreteRepository{db: tx, conn: r.conn}
}

func (r *concreteRepository) Create(ctx context.Context, enrollment Enrollment) (Enrollment, error) {
	if enrollment.CreatedAt.IsZero() {
		enrollment.CreatedAt = time.Now().UTC()
	}
	if enrollment.UpdatedAt.IsZero() {
		enrollment.UpdatedAt = enrollment.CreatedAt
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO course_members (course_id, user_id, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (course_id, user_id) DO UPDATE SET role = EXCLUDED.role, updated_at = EXCLUDED.updated_at`,
		enrollment.CourseID,
		enrollment.UserID,
		string(enrollment.Role),
		enrollment.CreatedAt,
		enrollment.UpdatedAt,
	)
	if err != nil {
		return Enrollment{}, err
	}
	return enrollment, nil
}

func (r *concreteRepository) Delete(ctx context.Context, courseID, userID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM course_members WHERE course_id = $1 AND user_id = $2`,
		courseID,
		userID,
	)
	return err
}

func (r *concreteRepository) Get(ctx context.Context, courseID, userID string) (Enrollment, error) {
	var e Enrollment
	err := r.db.QueryRowContext(
		ctx,
		`SELECT course_id, user_id, role, created_at, updated_at FROM course_members WHERE course_id = $1 AND user_id = $2`,
		courseID,
		userID,
	).Scan(
		&e.CourseID,
		&e.UserID,
		&e.Role,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Enrollment{}, ErrNotFound
		}
		return Enrollment{}, err
	}
	return e, nil
}

func (r *concreteRepository) ListByCourse(ctx context.Context, courseID string) ([]Enrollment, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT course_id, user_id, role, created_at, updated_at FROM course_members WHERE course_id = $1`,
		courseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(&e.CourseID, &e.UserID, &e.Role, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, nil
}
