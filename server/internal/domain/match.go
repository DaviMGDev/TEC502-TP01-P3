package domain

import (
	"errors"
	"sync"
)

const RoundsToWin = 2 // Best of 3

type MatchInterface interface {
	GetID() string
	GetPlayers() []UserInterface
	GetScores() map[string]int
	AddPlayer(player UserInterface) error
	RemovePlayer(playerID string) error
	MakeMove(playerID string, move CardInterface) error
	Surrender(playerID string) error
	GetWinner() (string, error)
}

type Match struct {
	ID      string                     `json:"id"`
	Players []UserInterface            `json:"players"`
	Moves   []map[string]CardInterface `json:"moves"`
	Scores  map[string]int             `json:"scores"`
	Winner  string                     `json:"winner"`
	mu      sync.RWMutex               `json:"-"`
}

func (m *Match) GetID() string {
	return m.ID
}

func (m *Match) GetPlayers() []UserInterface {
	return m.Players
}

func (m *Match) GetScores() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Return a copy to prevent race conditions on the caller's side
	scoresCopy := make(map[string]int)
	for k, v := range m.Scores {
		scoresCopy[k] = v
	}
	return scoresCopy
}

func (m *Match) AddPlayer(player UserInterface) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.Players) >= 2 {
		return errors.New("match is full")
	}
	m.Players = append(m.Players, player)
	m.Scores[player.GetID()] = 0 // Initialize score
	return nil
}

func (m *Match) RemovePlayer(playerID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, player := range m.Players {
		if player.GetID() == playerID {
			m.Players = append(m.Players[:i], m.Players[i+1:]...)
			delete(m.Scores, playerID)
			return nil
		}
	}
	return errors.New("player not found in match")
}

func (m *Match) MakeMove(playerID string, move CardInterface) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Winner != "" {
		return errors.New("match has already ended")
	}

	if len(m.Players) < 2 {
		return errors.New("not enough players to make a move")
	}

	if len(m.Moves) == 0 || len(m.Moves[len(m.Moves)-1]) == 2 {
		m.Moves = append(m.Moves, make(map[string]CardInterface))
	}

	currentRound := m.Moves[len(m.Moves)-1]
	if _, exists := currentRound[playerID]; exists {
		return errors.New("player has already made a move this round")
	}

	currentRound[playerID] = move
	if len(currentRound) == 2 {
		var player1ID, player2ID string
		var player1Move, player2Move CardInterface
		i := 0
		for pid, mv := range currentRound {
			if i == 0 {
				player1ID = pid
				player1Move = mv
			} else {
				player2ID = pid
				player2Move = mv
			}
			i++
		}

		result := player1Move.Against(player2Move)

		var roundWinnerID string
		if result == 1 {
			roundWinnerID = player1ID
		} else if result == -1 {
			roundWinnerID = player2ID
		}

		if roundWinnerID != "" {
			m.Scores[roundWinnerID]++
			if m.Scores[roundWinnerID] >= RoundsToWin {
				m.Winner = roundWinnerID
			}
		}
	}
	return nil
}

func (m *Match) Surrender(playerID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Winner != "" {
		return errors.New("match has already ended")
	}

	if len(m.Players) < 2 {
		return errors.New("not enough players in the match to determine a winner")
	}

	var winnerID string
	found := false
	for _, p := range m.Players {
		if p.GetID() == playerID {
			found = true
		} else {
			winnerID = p.GetID()
		}
	}

	if !found {
		return errors.New("surrendering player not found in match")
	}

	m.Winner = winnerID
	return nil
}

func (m *Match) GetWinner() (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.Winner == "" {
		return "", errors.New("no winner yet")
	}
	return m.Winner, nil
}

