package services

import (
	"cod-server/internal/domain"
)

type UserServiceInterface interface {
	Register(username, password string) error
	Login(username, password string) (*domain.UserInterface, error)
}

type CardsServiceInterface interface {
	GetCards(userID string) ([]domain.CardInterface, error)
	BuyPack(userID string) error
	OfferTrade(fromUserID, toUserID, cardID string) error
	AcceptTrade(fromUserID, toUserID, cardID string) error
}

type MatchServiceInterface interface {
	StartMatch(userID string) (domain.MatchInterface, error)
	JoinMatch(userID, matchID string) error
	SurrenderMatch(userID, matchID string) error
	MakeMove(userID, matchID string, cardID string) error
}

