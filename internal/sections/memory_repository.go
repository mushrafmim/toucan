package sections

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewMemoryRepository() *Repository {
	return &Repository{db: nil}
}
