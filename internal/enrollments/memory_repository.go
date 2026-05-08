package enrollments

import (
	"context"
	"database/sql"
	"sync"
)

type memoryRepository struct {
	items map[string]map[string]Enrollment // courseID -> userID -> Enrollment
	mu    sync.RWMutex
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		items: make(map[string]map[string]Enrollment),
	}
}

func (r *memoryRepository) DB() *sql.DB {
	return nil
}

func (r *memoryRepository) WithTx(tx *sql.Tx) Repository {
	return r
}

func (r *memoryRepository) Create(ctx context.Context, e Enrollment) (Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[e.CourseID]; !ok {
		r.items[e.CourseID] = make(map[string]Enrollment)
	}
	r.items[e.CourseID][e.UserID] = e
	return e, nil
}

func (r *memoryRepository) Delete(ctx context.Context, courseID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[courseID]; ok {
		delete(r.items[courseID], userID)
	}
	return nil
}

func (r *memoryRepository) Get(ctx context.Context, courseID, userID string) (Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if userMap, ok := r.items[courseID]; ok {
		if e, ok := userMap[userID]; ok {
			return e, nil
		}
	}
	return Enrollment{}, ErrNotFound
}

func (r *memoryRepository) ListByCourse(ctx context.Context, courseID string) ([]Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Enrollment
	if userMap, ok := r.items[courseID]; ok {
		for _, e := range userMap {
			result = append(result, e)
		}
	}
	return result, nil
}
