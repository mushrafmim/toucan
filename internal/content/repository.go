package content

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"toucan/internal/database"
)

type Repository interface {
	List(filter ListFilter) ListResult
	Get(id string) (Item, error)
	Create(item Item) (Item, error)
	Update(item Item) (Item, error)
	Delete(id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(filter ListFilter) ListResult {
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
	if filter.SectionID != "" {
		where = " WHERE section_id = $1"
		args = append(args, filter.SectionID)
	}

	countQuery := "SELECT COUNT(*) FROM content_items" + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return ListResult{Items: []Item{}, Page: page, PageSize: pageSize}
	}

	limitPos := len(args) + 1
	offsetPos := len(args) + 2
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	query := fmt.Sprintf(
		`SELECT id, section_id, title, summary, type, position, source_url, body, metadata, created_at, updated_at
		 FROM content_items%s
		 ORDER BY position ASC, created_at ASC, id ASC
		 LIMIT $%d OFFSET $%d`,
		where,
		limitPos,
		offsetPos,
	)

	rows, err := r.db.QueryContext(ctx, query, listArgs...)
	if err != nil {
		return ListResult{Items: []Item{}, Page: page, PageSize: pageSize, Total: total}
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		item, scanErr := scanContentItem(rows)
		if scanErr != nil {
			return ListResult{Items: []Item{}, Page: page, PageSize: pageSize, Total: total}
		}
		items = append(items, item)
	}

	return ListResult{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

func (r *repository) Get(id string) (Item, error) {
	row := r.db.QueryRowContext(
		context.Background(),
		`SELECT id, section_id, title, summary, type, position, source_url, body, metadata, created_at, updated_at
		 FROM content_items WHERE id = $1`,
		id,
	)

	item, err := scanContentItem(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return Item{}, ErrNotFound
		}
		return Item{}, err
	}
	return item, nil
}

func (r *repository) Create(item Item) (Item, error) {
	if item.ID == "" {
		item.ID = database.NewID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = item.CreatedAt
	}

	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return Item{}, err
	}

	_, err = r.db.ExecContext(
		context.Background(),
		`INSERT INTO content_items (
			id, section_id, title, summary, type, position, source_url, body, metadata, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		item.ID,
		item.SectionID,
		item.Title,
		item.Summary,
		string(item.Type),
		item.Position,
		database.NullString(item.SourceURL),
		database.NullString(item.Body),
		metadataJSON,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func (r *repository) Update(item Item) (Item, error) {
	item.UpdatedAt = time.Now().UTC()
	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return Item{}, err
	}

	result, err := r.db.ExecContext(
		context.Background(),
		`UPDATE content_items
		 SET title = $2, summary = $3, type = $4, position = $5, source_url = $6, body = $7, metadata = $8, updated_at = $9
		 WHERE id = $1`,
		item.ID,
		item.Title,
		item.Summary,
		string(item.Type),
		item.Position,
		database.NullString(item.SourceURL),
		database.NullString(item.Body),
		metadataJSON,
		item.UpdatedAt,
	)
	if err != nil {
		return Item{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return Item{}, err
	}
	if rows == 0 {
		return Item{}, ErrNotFound
	}
	return item, nil
}

func (r *repository) Delete(id string) error {
	result, err := r.db.ExecContext(context.Background(), `DELETE FROM content_items WHERE id = $1`, id)
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

type contentScanner interface {
	Scan(dest ...any) error
}

func scanContentItem(scanner contentScanner) (Item, error) {
	var item Item
	var sourceURL sql.NullString
	var body sql.NullString
	var metadataJSON []byte

	err := scanner.Scan(
		&item.ID,
		&item.SectionID,
		&item.Title,
		&item.Summary,
		&item.Type,
		&item.Position,
		&sourceURL,
		&body,
		&metadataJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return Item{}, err
	}
	if sourceURL.Valid {
		item.SourceURL = sourceURL.String
	}
	if body.Valid {
		item.Body = body.String
	}
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			return Item{}, err
		}
	}
	return item, nil
}
