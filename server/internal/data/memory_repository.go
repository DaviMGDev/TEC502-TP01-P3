package data

import (
	"errors"
	"sync"
)

// MemoryRepository é uma implementação em memória da interface Repository genérica.
type MemoryRepository[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

// NewMemoryRepository constrói um MemoryRepository com mapa vazio.
func NewMemoryRepository[T any]() Repository[T] {
	return &MemoryRepository[T]{
		data: make(map[string]T),
	}
}

// Create adiciona nova entidade; retorna erro se o id já existe.
func (r *MemoryRepository[T]) Create(id string, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; exists {
		return errors.New("entity with this id already exists")
	}

	r.data[id] = entity
	return nil
}

// Read obtém entidade pelo id; retorna erro se não encontrada.
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

// Update substitui entidade existente pelo id; retorna erro se id não existe.
func (r *MemoryRepository[T]) Update(id string, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return errors.New("entity not found")
	}

	r.data[id] = entity
	return nil
}

// Delete remove entidade pelo id; retorna erro se id não existe.
func (r *MemoryRepository[T]) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return errors.New("entity not found")
	}

	delete(r.data, id)
	return nil
}

// List retorna todas as entidades armazenadas.
func (r *MemoryRepository[T]) List() ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]T, 0, len(r.data))
	for _, entity := range r.data {
		list = append(list, entity)
	}

	return list, nil
}

// ListBy retorna entidades que casam com o filtro fornecido.
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
