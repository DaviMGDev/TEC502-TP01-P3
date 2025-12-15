package services

import (
	"cod-server/internal/domain"
)

// UserServiceInterface define métodos para operações de gerenciamento de contas de usuário.
type UserServiceInterface interface {
	// Register cria uma nova conta de usuário com as credenciais fornecidas.
	Register(username, password string) error
	// Login autentica um usuário e retorna seu objeto de domínio se bem-sucedido.
	Login(username, password string) (*domain.UserInterface, error)
}

// CardsServiceInterface define métodos para operações de propriedade e troca de cartas.
type CardsServiceInterface interface {
	// GetCards recupera todas as cartas pertencentes a um usuário específico.
	GetCards(userID string) ([]domain.CardInterface, error)
	// BuyPack permite que um usuário compre um novo pacote de cartas.
	BuyPack(userID string) error
	// OfferTrade inicia uma troca de carta de um usuário para outro.
	OfferTrade(fromUserID, toUserID, cardID string) error
	// AcceptTrade completa uma troca de carta oferecida anteriormente.
	AcceptTrade(fromUserID, toUserID, cardID string) error
}

// MatchServiceInterface define métodos para gerenciamento de partidas do jogo.
type MatchServiceInterface interface {
	// StartMatch cria e inicia uma nova partida de jogo para um usuário.
	StartMatch(userID string) (domain.MatchInterface, error)
	// JoinMatch permite que um usuário participe de uma partida existente.
	JoinMatch(userID, matchID string) error
	// SurrenderMatch encerra a partida atual com uma derrota para o usuário.
	SurrenderMatch(userID, matchID string) error
	// MakeMove joga uma carta durante uma partida ativa.
	MakeMove(userID, matchID string, cardID string) error
}
