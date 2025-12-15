package data

import (
	"errors"
	"sync"
)

// MemoryRepository é uma implementação em memória da interface Repository
type MemoryRepository[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

// NewMemoryRepository cria uma nova instância de MemoryRepository
func NewMemoryRepository[T any]() Repository[T] {
	return &MemoryRepository[T]{
		data: make(map[string]T),
	}
}

// Create adiciona uma nova entidade ao repositório
func (r *MemoryRepository[T]) Create(id string, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; exists {
		return errors.New("entity with this id already exists")
	}

	r.data[id] = entity
	return nil
}

// Read recupera uma entidade pelo ID
func (r *MemoryRepository[T]) Read(id string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, exists := r.data[id]
	if !exists {
		var zero T
		return zero, errors.New("entity not found")
	}

	return entity, nil
}

// Update atualiza uma entidade existente
func (r *MemoryRepository[T]) Update(id string, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return errors.New("entity not found")
	}

	r.data[id] = entity
	return nil
}

// Delete remove uma entidade do repositório
func (r *MemoryRepository[T]) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return errors.New("entity not found")
	}

	delete(r.data, id)
	return nil
}

// List retorna todas as entidades do repositório
func (r *MemoryRepository[T]) List() ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]T, 0, len(r.data))
	for _, entity := range r.data {
		list = append(list, entity)
	}

	return list, nil
}

// ListBy retorna uma lista de entidades que correspondem a um filtro
func (r *MemoryRepository[T]) ListBy(filter func(T) bool) ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]T, 0)
	for _, entity := range r.data {
		if filter(entity) {
			list = append(list, entity)
		}
	}

	return list, nil
}
