package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"errors"

	uuid "github.com/google/uuid"
)

type MatchService struct {
	matchRepo data.Repository[domain.MatchInterface]
	cardsRepo data.Repository[domain.CardInterface]
	usersRepo data.Repository[domain.UserInterface]
}

func NewMatchService(matchRepo data.Repository[domain.MatchInterface], cardsRepo data.Repository[domain.CardInterface], usersRepo data.Repository[domain.UserInterface]) MatchServiceInterface {
	return &MatchService{
		matchRepo: matchRepo,
		cardsRepo: cardsRepo,
		usersRepo: usersRepo,
	}
}

func (ms *MatchService) StartMatch(userID string) (domain.MatchInterface, error) {
	user, err := ms.usersRepo.Read(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	id := uuid.New().String()
	match := &domain.Match{
		ID:      id,
		Players: []domain.UserInterface{},
		Moves:   []map[string]domain.CardInterface{},
		Scores:  make(map[string]int),
		Winner:  "",
	}
	match.AddPlayer(user)

	err = ms.matchRepo.Create(id, match)
	if err != nil {
		return nil, err
	}
	var imatch domain.MatchInterface = match
	return imatch, nil
}

func (ms *MatchService) JoinMatch(userID, matchID string) error {
	match_raw, err := ms.matchRepo.Read(matchID)
	if err != nil {
		return err
	}

	user, err := ms.usersRepo.Read(userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = match_raw.AddPlayer(user)
	if err != nil {
		return err
	}

	return ms.matchRepo.Update(matchID, match_raw)
}

func (ms *MatchService) SurrenderMatch(userID, matchID string) error {
	match_raw, err := ms.matchRepo.Read(matchID)
	if err != nil {
		return err
	}

	err = match_raw.Surrender(userID)
	if err != nil {
		return err
	}

	return ms.matchRepo.Update(matchID, match_raw)
}

func (ms *MatchService) MakeMove(userID, matchID string, cardID string) error {
	match_raw, err := ms.matchRepo.Read(matchID)
	if err != nil {
		return err
	}

	card, err := ms.cardsRepo.Read(cardID)
	if err != nil {
		return err
	}

	// Domain object should check for card ownership
	if card.GetOwnerID() != userID {
		return errors.New("player does not own this card")
	}

	err = match_raw.MakeMove(userID, card)
	if err != nil {
		return err
	}
	return ms.matchRepo.Update(matchID, match_raw)
}

