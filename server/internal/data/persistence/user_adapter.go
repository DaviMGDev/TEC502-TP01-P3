package persistence

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
)

// UserRepoAdapter adapts a SQL-backed UserRepository to the generic data.Repository[domain.UserInterface].
type UserRepoAdapter struct {
	repo UserRepository
}

func NewUserRepoAdapter(repo UserRepository) data.Repository[domain.UserInterface] {
	return &UserRepoAdapter{repo: repo}
}

func (a *UserRepoAdapter) Create(id string, entity domain.UserInterface) error {
	// Convert interface to *domain.User when possible using type assertion.
	// Otherwise, construct a minimal *domain.User from interface getters.
	user, ok := entity.(*domain.User)
	if !ok {
		// Build a struct from interface methods when direct casting is not possible.
		return a.repo.Create(id, &domain.User{
			ID:       entity.GetID(),
			Username: entity.GetUsername(),
			// Precisamos de uma forma de extrair o Password e Cards
			// Isso pode exigir métodos adicionais na interface ou uma abordagem diferente
			Password: "", // Isso não é ideal
		})
	}
	return a.repo.Create(id, user)
}

func (a *UserRepoAdapter) Read(id string) (domain.UserInterface, error) {
	user, err := a.repo.Read(id)
	if err != nil {
		return nil, err
	}
	// O ponteiro *domain.User implementa UserInterface implicitamente
	return user, nil
}

func (a *UserRepoAdapter) Update(id string, entity domain.UserInterface) error {
	user, ok := entity.(*domain.User)
	if !ok {
		return a.repo.Update(id, &domain.User{
			ID:       entity.GetID(),
			Username: entity.GetUsername(),
			// Similar conversion logic applies when updating from the interface type.
			Password: "", // Isso não é ideal
		})
	}
	return a.repo.Update(id, user)
}

func (a *UserRepoAdapter) Delete(id string) error {
	return a.repo.Delete(id)
}

func (a *UserRepoAdapter) List() ([]domain.UserInterface, error) {
	users, err := a.repo.List()
	if err != nil {
		return nil, err
	}

	interfaces := make([]domain.UserInterface, len(users))
	for i, user := range users {
		interfaces[i] = user // *domain.User implementa UserInterface
	}
	return interfaces, nil
}

func (a *UserRepoAdapter) ListBy(filter func(domain.UserInterface) bool) ([]domain.UserInterface, error) {
	users, err := a.repo.List()
	if err != nil {
		return nil, err
	}

	var result []domain.UserInterface
	for _, user := range users {
		if filter(user) { // user já implementa UserInterface
			result = append(result, user)
		}
	}
	return result, nil
}
