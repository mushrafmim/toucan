package courses

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"toucan/internal/database"
)

type Repository interface {
	List(filter ListFilter) ListResult
	Get(id string) (Course, error)
	Create(course Course) (Course, error)
	Update(course Course) (Course, error)
	Delete(id string) error
	DB() *sql.DB
	WithTx(tx *sql.Tx) Repository
}

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

func (r *concreteRepository) List(filter ListFilter) ListResult {
	page := filter.Page
	if page <= 0 {
		page = 1
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	ctx := context.Background()
	where := make([]string, 0, 2)
	args := make([]any, 0, 4)
	argIndex := 1

	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(filter.Status))
		argIndex++
	}

	if query := strings.TrimSpace(filter.Query); query != "" {
		pattern := "%" + strings.ToLower(query) + "%"
		where = append(where, fmt.Sprintf(`(
			LOWER(title) LIKE $%d OR
			LOWER(summary) LIKE $%d OR
			LOWER(description) LIKE $%d OR
			LOWER(category) LIKE $%d OR
			LOWER(slug) LIKE $%d
		)`, argIndex, argIndex, argIndex, argIndex, argIndex))
		args = append(args, pattern)
		argIndex++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM courses" + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return ListResult{Items: []Course{}, Page: page, PageSize: pageSize}
	}

	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	query := `
		SELECT id, title, slug, summary, description, category, level, tags, status, creator_id, created_at, updated_at, published_at
		FROM courses` + whereClause + `
		ORDER BY created_at DESC, id ASC
		LIMIT $` + fmt.Sprint(argIndex) + ` OFFSET $` + fmt.Sprint(argIndex+1)

	rows, err := r.db.QueryContext(ctx, query, listArgs...)
	if err != nil {
		return ListResult{Items: []Course{}, Page: page, PageSize: pageSize, Total: total}
	}
	defer rows.Close()

	items := make([]Course, 0)
	for rows.Next() {
		course, scanErr := scanCourse(rows)
		if scanErr != nil {
			return ListResult{Items: []Course{}, Page: page, PageSize: pageSize, Total: total}
		}
		items = append(items, course)
	}

	return ListResult{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

func (r *concreteRepository) Get(id string) (Course, error) {
	row := r.db.QueryRowContext(
		context.Background(),
		`SELECT id, title, slug, summary, description, category, level, tags, status, creator_id, created_at, updated_at, published_at
		 FROM courses WHERE id = $1`,
		id,
	)

	course, err := scanCourse(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return Course{}, ErrNotFound
		}
		return Course{}, err
	}
	return course, nil
}

func (r *concreteRepository) Create(course Course) (Course, error) {
	if course.ID == "" {
		course.ID = database.NewID()
	}
	if course.CreatedAt.IsZero() {
		course.CreatedAt = time.Now().UTC()
	}
	if course.UpdatedAt.IsZero() {
		course.UpdatedAt = course.CreatedAt
	}

	tagsJSON, err := json.Marshal(course.Tags)
	if err != nil {
		return Course{}, err
	}

	_, err = r.db.ExecContext(
		context.Background(),
		`INSERT INTO courses (
			id, title, slug, summary, description, category, level, tags, status, creator_id, created_at, updated_at, published_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		course.ID,
		course.Title,
		course.Slug,
		course.Summary,
		course.Description,
		course.Category,
		string(course.Level),
		tagsJSON,
		string(course.Status),
		course.CreatorID,
		course.CreatedAt,
		course.UpdatedAt,
		database.NullTime(course.PublishedAt),
	)
	if err != nil {
		return Course{}, err
	}
	return course, nil
}

func (r *concreteRepository) Update(course Course) (Course, error) {
	course.UpdatedAt = time.Now().UTC()
	tagsJSON, err := json.Marshal(course.Tags)
	if err != nil {
		return Course{}, err
	}

	result, err := r.db.ExecContext(
		context.Background(),
		`UPDATE courses
		 SET title = $2, slug = $3, summary = $4, description = $5, category = $6, level = $7,
		     tags = $8, status = $9, updated_at = $10, published_at = $11
		 WHERE id = $1`,
		course.ID,
		course.Title,
		course.Slug,
		course.Summary,
		course.Description,
		course.Category,
		string(course.Level),
		tagsJSON,
		string(course.Status),
		course.UpdatedAt,
		database.NullTime(course.PublishedAt),
	)
	if err != nil {
		return Course{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return Course{}, err
	}
	if rows == 0 {
		return Course{}, ErrNotFound
	}

	return course, nil
}

func (r *concreteRepository) Delete(id string) error {
	result, err := r.db.ExecContext(context.Background(), `DELETE FROM courses WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

type courseScanner interface {
	Scan(dest ...any) error
}

func scanCourse(scanner courseScanner) (Course, error) {
	var course Course
	var tagsJSON []byte
	var publishedAt sql.NullTime

	err := scanner.Scan(
		&course.ID,
		&course.Title,
		&course.Slug,
		&course.Summary,
		&course.Description,
		&course.Category,
		&course.Level,
		&tagsJSON,
		&course.Status,
		&course.CreatorID,
		&course.CreatedAt,
		&course.UpdatedAt,
		&publishedAt,
	)
	if err != nil {
		return Course{}, err
	}
	if len(tagsJSON) > 0 {
		if err := json.Unmarshal(tagsJSON, &course.Tags); err != nil {
			return Course{}, err
		}
	}
	if publishedAt.Valid {
		course.PublishedAt = publishedAt.Time
	}
	return course, nil
}
