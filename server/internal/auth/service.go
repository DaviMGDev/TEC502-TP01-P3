package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("cod-server-secret-key-change-in-production") // Em produção, isso deve vir de variáveis de ambiente

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	secret []byte
}

func NewAuthService(secret string) *AuthService {
	if secret == "" {
		secret = string(jwtSecret)
	}
	return &AuthService{secret: []byte(secret)}
}

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