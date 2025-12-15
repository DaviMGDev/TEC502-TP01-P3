package auth

// AuthService fornece utilitários de autenticação baseados em JWT para emitir e validar tokens.
// Em produção, o segredo deve ser injetado via variáveis de ambiente ou um cofre seguro.

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("cod-server-secret-key-change-in-production") // Definir via env em produção

// Claims estende claims registrados do JWT com campos específicos da aplicação.
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	secret []byte
}

// NewAuthService constrói um AuthService, usando um segredo interno padrão se vazio.
func NewAuthService(secret string) *AuthService {
	if secret == "" {
		secret = string(jwtSecret)
	}
	return &AuthService{secret: []byte(secret)}
}

// GenerateToken cria um JWT assinado para o usuário fornecido com expiração de 24h.
func (s *AuthService) GenerateToken(userID, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken analisa e valida um JWT, retornando claims quando válido.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
