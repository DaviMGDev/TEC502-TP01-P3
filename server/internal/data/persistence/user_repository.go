package persistence

// SqlUserRepository implementa UserRepository usando um banco SQLite.

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"cod-server/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type SqlUserRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) UserRepository {
	// Garante que a tabela users exista na inicialização.
	// Create insere um novo usuário com cartas serializadas.
	// Read busca um usuário por id; retorna erro se não encontrar. Desserialização de cartas está simplificada.
	// Update modifica username, password e cartas serializadas para o id fornecido.
	// Delete remove um usuário; retorna erro se nenhuma linha for afetada.
	// List recupera todos os usuários do banco.
	// ListBy filtra usuários em memória usando o predicado fornecido.
	query := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		cards TEXT
	)`
	_, err := db.Exec(query)
	if err != nil {
		panic(fmt.Sprintf("Falha ao criar tabela users: %v", err))
	}

	return &SqlUserRepository{db: db}
}

func (r *SqlUserRepository) Create(id string, user *domain.User) error {
	cardsJSON, err := json.Marshal(user.Cards)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT INTO users (id, username, password, cards) VALUES (?, ?, ?, ?)",
		id, user.Username, user.Password, string(cardsJSON))
	return err
}

func (r *SqlUserRepository) Read(id string) (*domain.User, error) {
	var user domain.User
	var cardsJSON string

	err := r.db.QueryRow("SELECT id, username, password, cards FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Username, &user.Password, &cardsJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Obs.: o campo Cards deveria ser tratado adequadamente - isto é uma simplificação
	// Seria necessário implementar serialização adequada para PackInterface

	return &user, nil
}

func (r *SqlUserRepository) Update(id string, user *domain.User) error {
	cardsJSON, err := json.Marshal(user.Cards)
	if err != nil {
		return err
	}

	result, err := r.db.Exec("UPDATE users SET username = ?, password = ?, cards = ? WHERE id = ?",
		user.Username, user.Password, string(cardsJSON), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *SqlUserRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *SqlUserRepository) List() ([]*domain.User, error) {
	rows, err := r.db.Query("SELECT id, username, password, cards FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var cardsJSON string

		err := rows.Scan(&user.ID, &user.Username, &user.Password, &cardsJSON)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *SqlUserRepository) ListBy(filter func(*domain.User) bool) ([]*domain.User, error) {
	allUsers, err := r.List()
	if err != nil {
		return nil, err
	}

	var filteredUsers []*domain.User
	for _, user := range allUsers {
		if filter(user) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers, nil
}
