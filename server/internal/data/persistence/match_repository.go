package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"cod-server/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type SqlMatchRepository struct {
	db *sql.DB
}

func NewSqlMatchRepository(db *sql.DB) MatchRepository {
	// Create table if not exists
	query := `CREATE TABLE IF NOT EXISTS matches (
		id TEXT PRIMARY KEY,
		players TEXT,
		moves TEXT,
		scores TEXT,
		winner TEXT
	)`
	_, err := db.Exec(query)
	if err != nil {
		panic(fmt.Sprintf("Failed to create matches table: %v", err))
	}

	return &SqlMatchRepository{db: db}
}

func (r *SqlMatchRepository) Create(id string, match *domain.Match) error {
	playersJSON, err := json.Marshal(match.Players)
	if err != nil {
		return err
	}
	
	movesJSON, err := json.Marshal(match.Moves)
	if err != nil {
		return err
	}
	
	scoresJSON, err := json.Marshal(match.Scores)
	if err != nil {
		return err
	}
	
	_, err = r.db.Exec("INSERT INTO matches (id, players, moves, scores, winner) VALUES (?, ?, ?, ?, ?)", 
		id, string(playersJSON), string(movesJSON), string(scoresJSON), match.Winner)
	return err
}

func (r *SqlMatchRepository) Read(id string) (*domain.Match, error) {
	var match domain.Match
	var playersJSON, movesJSON, scoresJSON string
	
	err := r.db.QueryRow("SELECT id, players, moves, scores, winner FROM matches WHERE id = ?", id).
		Scan(&match.ID, &playersJSON, &movesJSON, &scoresJSON, &match.Winner)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("match not found")
		}
		return nil, err
	}
	
	// Note: Deserializing complex nested structures like Players, Moves would require 
	// more complex handling
	
	return &match, nil
}

func (r *SqlMatchRepository) Update(id string, match *domain.Match) error {
	playersJSON, err := json.Marshal(match.Players)
	if err != nil {
		return err
	}
	
	movesJSON, err := json.Marshal(match.Moves)
	if err != nil {
		return err
	}
	
	scoresJSON, err := json.Marshal(match.Scores)
	if err != nil {
		return err
	}
	
	result, err := r.db.Exec("UPDATE matches SET players = ?, moves = ?, scores = ?, winner = ? WHERE id = ?", 
		string(playersJSON), string(movesJSON), string(scoresJSON), match.Winner, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("match not found")
	}
	
	return nil
}

func (r *SqlMatchRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM matches WHERE id = ?", id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("match not found")
	}
	
	return nil
}

func (r *SqlMatchRepository) List() ([]*domain.Match, error) {
	rows, err := r.db.Query("SELECT id, players, moves, scores, winner FROM matches")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var matches []*domain.Match
	for rows.Next() {
		var match domain.Match
		var playersJSON, movesJSON, scoresJSON string
		
		err := rows.Scan(&match.ID, &playersJSON, &movesJSON, &scoresJSON, &match.Winner)
		if err != nil {
			return nil, err
		}
		
		matches = append(matches, &match)
	}
	
	return matches, nil
}

func (r *SqlMatchRepository) ListBy(filter func(*domain.Match) bool) ([]*domain.Match, error) {
	allMatches, err := r.List()
	if err != nil {
		return nil, err
	}
	
	var filteredMatches []*domain.Match
	for _, match := range allMatches {
		if filter(match) {
			filteredMatches = append(filteredMatches, match)
		}
	}
	
	return filteredMatches, nil
}