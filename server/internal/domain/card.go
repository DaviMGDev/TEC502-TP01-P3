package domain

import (
	"errors"
)

var (
	ErrCardNotOwnedByUser = errors.New("card is not owned by user")
)

type CardInterface interface {
	GetID() string
	GetOwnerID() string
	GetType() string
	Against(opponent CardInterface) int
}

type PackInterface interface {
	GetID() string
	GetCards() []CardInterface
	AddCard(card CardInterface)
	DrawCard(index int) (CardInterface, error)
}

type Card struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Type    string `json:"type"`
}

type Pack struct {
	ID    string          `json:"id"`
	Cards []CardInterface `json:"cards"`
}

func (c *Card) GetID() string {
	return c.ID
}

func (c *Card) GetOwnerID() string {
	return c.OwnerID
}

func (c *Card) GetType() string {
	return c.Type
}

func (c *Card) Against(opponent CardInterface) int {
	if c.Type == opponent.GetType() {
		return 0 // Draw
	}
	if (c.Type == "rock" && opponent.GetType() == "scissors") ||
		(c.Type == "scissors" && opponent.GetType() == "paper") ||
		(c.Type == "paper" && opponent.GetType() == "rock") {
		return 1 // Win
	}
	return -1 // Lose
}

func (p *Pack) GetID() string {
	return p.ID
}

func (p *Pack) GetCards() []CardInterface {
	return p.Cards
}

func (p *Pack) AddCard(card CardInterface) {
	p.Cards = append(p.Cards, card)
}

func (p *Pack) DrawCard(index int) (CardInterface, error) {
	if index < 0 || index >= len(p.Cards) {
		return nil, errors.New("index out of range")
	}
	card := p.Cards[index]
	p.Cards = append(p.Cards[:index], p.Cards[index+1:]...)
	return card, nil
}

