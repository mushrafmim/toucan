package database

import (
	"database/sql"
	"fmt"
)

func EnsureSchema(db *sql.DB) error {
	statements := []string{
		`
		CREATE TABLE IF NOT EXISTS courses (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			slug TEXT NOT NULL,
			summary TEXT NOT NULL,
			description TEXT NOT NULL,
			category TEXT NOT NULL,
			level TEXT NOT NULL,
			tags JSONB NOT NULL DEFAULT '[]'::jsonb,
			status TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			published_at TIMESTAMPTZ NULL
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS sections (
			id TEXT PRIMARY KEY,
			course_id TEXT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			summary TEXT NOT NULL,
			position INTEGER NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS content_items (
			id TEXT PRIMARY KEY,
			section_id TEXT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			summary TEXT NOT NULL,
			type TEXT NOT NULL,
			position INTEGER NOT NULL,
			source_url TEXT NULL,
			body TEXT NULL,
			metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)
		`,
		`CREATE INDEX IF NOT EXISTS idx_courses_status_created_at ON courses(status, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_sections_course_id_position ON sections(course_id, position, created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_content_items_section_id_position ON content_items(section_id, position, created_at)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("ensure postgres schema: %w", err)
		}
	}

	return nil
}
