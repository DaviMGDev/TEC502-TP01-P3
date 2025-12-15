package persistence

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
)

// MatchRepoAdapter adapta um MatchRepository baseado em SQL para data.Repository[domain.MatchInterface].
type MatchRepoAdapter struct {
	repo MatchRepository
}

func NewMatchRepoAdapter(repo MatchRepository) data.Repository[domain.MatchInterface] {
	return &MatchRepoAdapter{repo: repo}
}

func (a *MatchRepoAdapter) Create(id string, entity domain.MatchInterface) error {
	match, ok := entity.(*domain.Match)
	if !ok {
		return a.repo.Create(id, &domain.Match{
			ID:      entity.GetID(),
			Players: entity.GetPlayers(),
			// Obs.: Campos complexos como Moves e Scores são omitidos ao construir via métodos de interface.
			// Uma implementação completa exigiria getters adicionais na interface.
		})
	}
	return a.repo.Create(id, match)
}

func (a *MatchRepoAdapter) Read(id string) (domain.MatchInterface, error) {
	match, err := a.repo.Read(id)
	if err != nil {
		return nil, err
	}
	return match, nil // *domain.Match implementa MatchInterface
}

func (a *MatchRepoAdapter) Update(id string, entity domain.MatchInterface) error {
	match, ok := entity.(*domain.Match)
	if !ok {
		return a.repo.Update(id, &domain.Match{
			ID:      entity.GetID(),
			Players: entity.GetPlayers(),
			// Problema similar com campos complexos exigindo getters adicionais na interface.
		})
	}
	return a.repo.Update(id, match)
}

func (a *MatchRepoAdapter) Delete(id string) error {
	return a.repo.Delete(id)
}

func (a *MatchRepoAdapter) List() ([]domain.MatchInterface, error) {
	matches, err := a.repo.List()
	if err != nil {
		return nil, err
	}

	interfaces := make([]domain.MatchInterface, len(matches))
	for i, match := range matches {
		interfaces[i] = match
	}
	return interfaces, nil
}

func (a *MatchRepoAdapter) ListBy(filter func(domain.MatchInterface) bool) ([]domain.MatchInterface, error) {
	matches, err := a.repo.List()
	if err != nil {
		return nil, err
	}

	var result []domain.MatchInterface
	for _, match := range matches {
		if filter(match) {
			result = append(result, match)
		}
	}
	return result, nil
}
