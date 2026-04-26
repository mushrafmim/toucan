package sections

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"toucan/internal/database"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(filter ListFilter) ListResult {
	page := max(filter.Page, 1)
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	ctx := context.Background()
	args := []any{}
	where := ""
	if filter.CourseID != "" {
		where = " WHERE course_id = $1"
		args = append(args, filter.CourseID)
	}

	countQuery := "SELECT COUNT(*) FROM sections" + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return ListResult{Items: []Section{}, Page: page, PageSize: pageSize}
	}

	limitPos := len(args) + 1
	offsetPos := len(args) + 2
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	query := fmt.Sprintf(
		`SELECT id, course_id, title, summary, position, created_at, updated_at
		 FROM sections%s
		 ORDER BY position ASC, created_at ASC, id ASC
		 LIMIT $%d OFFSET $%d`,
		where,
		limitPos,
		offsetPos,
	)

	rows, err := r.db.QueryContext(ctx, query, listArgs...)
	if err != nil {
		return ListResult{Items: []Section{}, Page: page, PageSize: pageSize, Total: total}
	}
	defer rows.Close()

	items := make([]Section, 0)
	for rows.Next() {
		section, scanErr := scanSection(rows)
		if scanErr != nil {
			return ListResult{Items: []Section{}, Page: page, PageSize: pageSize, Total: total}
		}
		items = append(items, section)
	}

	return ListResult{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

func (r *Repository) Get(id string) (Section, error) {
	row := r.db.QueryRowContext(
		context.Background(),
		`SELECT id, course_id, title, summary, position, created_at, updated_at
		 FROM sections WHERE id = $1`,
		id,
	)

	section, err := scanSection(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return Section{}, ErrNotFound
		}
		return Section{}, err
	}
	return section, nil
}

func (r *Repository) Create(section Section) (Section, error) {
	if section.ID == "" {
		section.ID = database.NewID()
	}
	if section.CreatedAt.IsZero() {
		section.CreatedAt = time.Now().UTC()
	}
	if section.UpdatedAt.IsZero() {
		section.UpdatedAt = section.CreatedAt
	}

	_, err := r.db.ExecContext(
		context.Background(),
		`INSERT INTO sections (id, course_id, title, summary, position, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		section.ID,
		section.CourseID,
		section.Title,
		section.Summary,
		section.Position,
		section.CreatedAt,
		section.UpdatedAt,
	)
	if err != nil {
		return Section{}, err
	}
	return section, nil
}

func (r *Repository) Update(section Section) (Section, error) {
	section.UpdatedAt = time.Now().UTC()

	result, err := r.db.ExecContext(
		context.Background(),
		`UPDATE sections
		 SET title = $2, summary = $3, position = $4, updated_at = $5
		 WHERE id = $1`,
		section.ID,
		section.Title,
		section.Summary,
		section.Position,
		section.UpdatedAt,
	)
	if err != nil {
		return Section{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return Section{}, err
	}
	if rows == 0 {
		return Section{}, ErrNotFound
	}
	return section, nil
}

func (r *Repository) Delete(id string) error {
	result, err := r.db.ExecContext(context.Background(), `DELETE FROM sections WHERE id = $1`, id)
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

type sectionScanner interface {
	Scan(dest ...any) error
}

func scanSection(scanner sectionScanner) (Section, error) {
	var section Section
	err := scanner.Scan(
		&section.ID,
		&section.CourseID,
		&section.Title,
		&section.Summary,
		&section.Position,
		&section.CreatedAt,
		&section.UpdatedAt,
	)
	if err != nil {
		return Section{}, err
	}
	return section, nil
}
