package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"

	"github.com/google/uuid"
)

type CardsService struct {
	cardsRepo data.Repository[domain.CardInterface]
	usersRepo data.Repository[domain.UserInterface]
}

func NewCardsService(cardsRepo data.Repository[domain.CardInterface], usersRepo data.Repository[domain.UserInterface]) CardsServiceInterface {
	return &CardsService{cardsRepo: cardsRepo, usersRepo: usersRepo}
}

func (cs *CardsService) GetCards(userID string) ([]domain.CardInterface, error) {
	cards, err := cs.cardsRepo.ListBy(func(c domain.CardInterface) bool {
		return c.GetOwnerID() == userID
	})
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cs *CardsService) BuyPack(userID string) error {
	// Verify user exists
	_, err := cs.usersRepo.Read(userID)
	if err != nil {
		return err
	}

	// Create 5 random cards for the user
	cardTypes := []string{"rock", "paper", "scissors", "rock", "paper"} // Example card types
	for i := 0; i < 5; i++ {
		cardID := uuid.New().String()
		card := &domain.Card{
			ID:      cardID,
			OwnerID: userID,
			Type:    cardTypes[i%len(cardTypes)],
		}
		err := cs.cardsRepo.Create(cardID, card)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs *CardsService) OfferTrade(fromUserID, toUserID, cardID string) error {
	// Verify that the from user exists
	_, err := cs.usersRepo.Read(fromUserID)
	if err != nil {
		return err
	}

	// Verify that the to user exists
	_, err = cs.usersRepo.Read(toUserID)
	if err != nil {
		return err
	}

	// Verify that the card exists and belongs to the from user
	card, err := cs.cardsRepo.Read(cardID)
	if err != nil {
		return err
	}

	if card.GetOwnerID() != fromUserID {
		return domain.ErrCardNotOwnedByUser
	}

	// Here we would typically store the trade offer in a separate repository
	// For now, we'll just validate that the trade request is valid
	// In a real implementation, we would create a trade offer entity and store it

	return nil
}

func (cs *CardsService) AcceptTrade(fromUserID, toUserID, cardID string) error {
	// Verify that the to user exists
	_, err := cs.usersRepo.Read(toUserID)
	if err != nil {
		return err
	}

	// Verify that the card exists and belongs to the from user
	card, err := cs.cardsRepo.Read(cardID)
	if err != nil {
		return err
	}

	if card.GetOwnerID() != fromUserID {
		return domain.ErrCardNotOwnedByUser
	}

	// Perform the trade - update the card's owner
	updatedCard := &domain.Card{
		ID:      card.GetID(),
		OwnerID: toUserID, // Give the card to the user who accepted the trade
		Type:    card.GetType(),
	}

	// Update the card in the repository
	err = cs.cardsRepo.Update(cardID, updatedCard)
	if err != nil {
		return err
	}

	return nil
}	

