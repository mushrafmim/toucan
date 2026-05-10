package courses

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"
)

type memoryRepository struct {
	items map[string]Course
	mu    sync.RWMutex
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		items: make(map[string]Course),
	}
}

func (r *memoryRepository) DB() *sql.DB {
	return nil
}

func (r *memoryRepository) WithTx(tx *sql.Tx) Repository {
	return r
}

func (r *memoryRepository) List(filter ListFilter) (ListResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []Course
	for _, item := range r.items {
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.Query != "" && !strings.Contains(strings.ToLower(item.Title), strings.ToLower(filter.Query)) {
			continue
		}
		// Note: Memory repo doesn't fully support UserID join for now,
		// but we can add it if needed for tests.
		items = append(items, item)
	}

	return ListResult{
		Items:    items,
		Page:     1,
		PageSize: 10,
		Total:    len(items),
	}, nil
}

func (r *memoryRepository) Get(id string) (Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return Course{}, ErrNotFound
	}
	return item, nil
}

func (r *memoryRepository) Create(course Course) (Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if course.ID == "" {
		course.ID = fmt.Sprintf("course-%d", len(r.items)+1)
	}
	course.CreatedAt = time.Now().UTC()
	course.UpdatedAt = course.CreatedAt
	r.items[course.ID] = course
	return course, nil
}

func (r *memoryRepository) Update(course Course) (Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[course.ID]; !ok {
		return Course{}, ErrNotFound
	}
	course.UpdatedAt = time.Now().UTC()
	r.items[course.ID] = course
	return course, nil
}

func (r *memoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.items, id)
	return nil
}
