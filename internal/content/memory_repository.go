package content

import (
	"sync"
)

type memoryRepository struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		items: make(map[string]Item),
	}
}

func (r *memoryRepository) List(filter ListFilter) ListResult {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []Item
	for _, item := range r.items {
		if filter.SectionID != "" && item.SectionID != filter.SectionID {
			continue
		}
		matches = append(matches, item)
	}

	return ListResult{
		Items:    matches,
		Total:    len(matches),
		Page:     1,
		PageSize: 100,
	}
}

func (r *memoryRepository) Get(id string) (Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return Item{}, ErrNotFound
	}
	return item, nil
}

func (r *memoryRepository) Create(item Item) (Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[item.ID] = item
	return item, nil
}

func (r *memoryRepository) Update(item Item) (Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[item.ID]; !ok {
		return Item{}, ErrNotFound
	}
	r.items[item.ID] = item
	return item, nil
}

func (r *memoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
