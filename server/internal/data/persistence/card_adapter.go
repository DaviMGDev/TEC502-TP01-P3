package persistence

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
)

// Adaptador para tornar SqlCardRepository compatível com data.Repository[domain.CardInterface]
type CardRepoAdapter struct {
	repo CardRepository
}

func NewCardRepoAdapter(repo CardRepository) data.Repository[domain.CardInterface] {
	return &CardRepoAdapter{repo: repo}
}

func (a *CardRepoAdapter) Create(id string, entity domain.CardInterface) error {
	card, ok := entity.(*domain.Card)
	if !ok {
		return a.repo.Create(id, &domain.Card{
			ID:      entity.GetID(),
			OwnerID: entity.GetOwnerID(),
			Type:    entity.GetType(),
		})
	}
	return a.repo.Create(id, card)
}

func (a *CardRepoAdapter) Read(id string) (domain.CardInterface, error) {
	card, err := a.repo.Read(id)
	if err != nil {
		return nil, err
	}
	return card, nil // *domain.Card implementa CardInterface
}

func (a *CardRepoAdapter) Update(id string, entity domain.CardInterface) error {
	card, ok := entity.(*domain.Card)
	if !ok {
		return a.repo.Update(id, &domain.Card{
			ID:      entity.GetID(),
			OwnerID: entity.GetOwnerID(),
			Type:    entity.GetType(),
		})
	}
	return a.repo.Update(id, card)
}

func (a *CardRepoAdapter) Delete(id string) error {
	return a.repo.Delete(id)
}

func (a *CardRepoAdapter) List() ([]domain.CardInterface, error) {
	cards, err := a.repo.List()
	if err != nil {
		return nil, err
	}
	
	interfaces := make([]domain.CardInterface, len(cards))
	for i, card := range cards {
		interfaces[i] = card
	}
	return interfaces, nil
}

func (a *CardRepoAdapter) ListBy(filter func(domain.CardInterface) bool) ([]domain.CardInterface, error) {
	cards, err := a.repo.List()
	if err != nil {
		return nil, err
	}
	
	var result []domain.CardInterface
	for _, card := range cards {
		if filter(card) {
			result = append(result, card)
		}
	}
	return result, nil
}