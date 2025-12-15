package services

import (
	"testing"
	"cod-server/internal/domain"
)

// MockUserRepository é um repositório mock para testes
type MockUserRepository struct {
	users map[string]domain.UserInterface
}

func (m *MockUserRepository) Create(id string, entity domain.UserInterface) error {
	if m.users == nil {
		m.users = make(map[string]domain.UserInterface)
	}
	m.users[id] = entity
	return nil
}

func (m *MockUserRepository) Read(id string) (domain.UserInterface, error) {
	if m.users == nil {
		m.users = make(map[string]domain.UserInterface)
	}
	user, exists := m.users[id]
	if !exists {
		return nil, nil // Simples implementação para testes
	}
	return user, nil
}

func (m *MockUserRepository) Update(id string, entity domain.UserInterface) error {
	if m.users == nil {
		m.users = make(map[string]domain.UserInterface)
	}
	m.users[id] = entity
	return nil
}

func (m *MockUserRepository) Delete(id string) error {
	if m.users == nil {
		m.users = make(map[string]domain.UserInterface)
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List() ([]domain.UserInterface, error) {
	var users []domain.UserInterface
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) ListBy(filter func(domain.UserInterface) bool) ([]domain.UserInterface, error) {
	var filteredUsers []domain.UserInterface
	for _, user := range m.users {
		if filter(user) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, nil
}

// MockCardRepository é um repositório mock para testes de cartas
type MockCardRepository struct {
	cards map[string]domain.CardInterface
}

func (m *MockCardRepository) Create(id string, entity domain.CardInterface) error {
	if m.cards == nil {
		m.cards = make(map[string]domain.CardInterface)
	}
	m.cards[id] = entity
	return nil
}

func (m *MockCardRepository) Read(id string) (domain.CardInterface, error) {
	if m.cards == nil {
		m.cards = make(map[string]domain.CardInterface)
	}
	card, exists := m.cards[id]
	if !exists {
		return nil, nil
	}
	return card, nil
}

func (m *MockCardRepository) Update(id string, entity domain.CardInterface) error {
	if m.cards == nil {
		m.cards = make(map[string]domain.CardInterface)
	}
	m.cards[id] = entity
	return nil
}

func (m *MockCardRepository) Delete(id string) error {
	if m.cards == nil {
		m.cards = make(map[string]domain.CardInterface)
	}
	delete(m.cards, id)
	return nil
}

func (m *MockCardRepository) List() ([]domain.CardInterface, error) {
	var cards []domain.CardInterface
	for _, card := range m.cards {
		cards = append(cards, card)
	}
	return cards, nil
}

func (m *MockCardRepository) ListBy(filter func(domain.CardInterface) bool) ([]domain.CardInterface, error) {
	var filteredCards []domain.CardInterface
	for _, card := range m.cards {
		if filter(card) {
			filteredCards = append(filteredCards, card)
		}
	}
	return filteredCards, nil
}

func TestUserService_Register(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := NewUserService(mockRepo)

	username := "testuser"
	password := "testpass"
	
	err := userService.Register(username, password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify user was created with hashed password
	users, err := mockRepo.List()
	if err != nil || len(users) == 0 {
		t.Error("User was not saved to repository")
		return
	}
	
	savedUser := users[0]
	if savedUser.GetUsername() != username {
		t.Errorf("Expected username %s, got %s", username, savedUser.GetUsername())
	}
	
	// Verify that the password hash is valid
	if !savedUser.CheckPassword(password) {
		t.Error("Password hash verification failed")
	}
}

func TestUserService_Login(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := NewUserService(mockRepo)

	// Create a user first
	username := "testuser"
	password := "testpass"
	userService.Register(username, password)
	
	// Try to login
	user, err := userService.Login(username, password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if user == nil {
		t.Error("Expected user, got nil")
	} else if (*user).GetUsername() != username {
		t.Errorf("Expected username %s, got %s", username, (*user).GetUsername())
	}
	
	// Try to login with wrong password
	user, err = userService.Login(username, "wrongpass")
	if err != nil && user != nil {
		t.Error("Expected nil user with wrong password")
	}
}

func TestCardsService_GetCards(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockUserRepo := &MockUserRepository{}
	cardsService := NewCardsService(mockCardRepo, mockUserRepo)

	// Create a user and some cards
	userID := "user123"
	user := &domain.User{
		ID:       userID,
		Username: "testuser",
		Password: "hashedpass",
	}

	// Inicializar o mapa se ainda não estiver inicializado
	if mockUserRepo.users == nil {
		mockUserRepo.users = make(map[string]domain.UserInterface)
	}
	mockUserRepo.users[userID] = user

	card1 := &domain.Card{
		ID:      "card1",
		OwnerID: userID,
		Type:    "rock",
	}
	card2 := &domain.Card{
		ID:      "card2",
		OwnerID: userID,
		Type:    "paper",
	}

	// Inicializar o mapa se ainda não estiver inicializado
	if mockCardRepo.cards == nil {
		mockCardRepo.cards = make(map[string]domain.CardInterface)
	}

	mockCardRepo.cards[card1.ID] = card1
	mockCardRepo.cards[card2.ID] = card2

	cards, err := cardsService.GetCards(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(cards) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(cards))
	}
	
	// Verify that all cards belong to the correct user
	for _, card := range cards {
		if card.GetOwnerID() != userID {
			t.Errorf("Expected card owner %s, got %s", userID, card.GetOwnerID())
		}
	}
}