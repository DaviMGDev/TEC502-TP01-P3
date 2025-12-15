package persistence

import (
	"database/sql"
	"fmt"

	"cod-server/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type SqlCardRepository struct {
	db *sql.DB
}

func NewSqlCardRepository(db *sql.DB) CardRepository {
	// Create table if not exists
	query := `CREATE TABLE IF NOT EXISTS cards (
		id TEXT PRIMARY KEY,
		owner_id TEXT NOT NULL,
		card_type TEXT NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		panic(fmt.Sprintf("Failed to create cards table: %v", err))
	}

	return &SqlCardRepository{db: db}
}

func (r *SqlCardRepository) Create(id string, card *domain.Card) error {
	_, err := r.db.Exec("INSERT INTO cards (id, owner_id, card_type) VALUES (?, ?, ?)", 
		id, card.OwnerID, card.Type)
	return err
}

func (r *SqlCardRepository) Read(id string) (*domain.Card, error) {
	var card domain.Card
	
	err := r.db.QueryRow("SELECT id, owner_id, card_type FROM cards WHERE id = ?", id).
		Scan(&card.ID, &card.OwnerID, &card.Type)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card not found")
		}
		return nil, err
	}
	
	return &card, nil
}

func (r *SqlCardRepository) Update(id string, card *domain.Card) error {
	result, err := r.db.Exec("UPDATE cards SET owner_id = ?, card_type = ? WHERE id = ?", 
		card.OwnerID, card.Type, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("card not found")
	}
	
	return nil
}

func (r *SqlCardRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM cards WHERE id = ?", id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("card not found")
	}
	
	return nil
}

func (r *SqlCardRepository) List() ([]*domain.Card, error) {
	rows, err := r.db.Query("SELECT id, owner_id, card_type FROM cards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var cards []*domain.Card
	for rows.Next() {
		var card domain.Card
		
		err := rows.Scan(&card.ID, &card.OwnerID, &card.Type)
		if err != nil {
			return nil, err
		}
		
		cards = append(cards, &card)
	}
	
	return cards, nil
}

func (r *SqlCardRepository) ListBy(filter func(*domain.Card) bool) ([]*domain.Card, error) {
	allCards, err := r.List()
	if err != nil {
		return nil, err
	}
	
	var filteredCards []*domain.Card
	for _, card := range allCards {
		if filter(card) {
			filteredCards = append(filteredCards, card)
		}
	}
	
	return filteredCards, nil
}