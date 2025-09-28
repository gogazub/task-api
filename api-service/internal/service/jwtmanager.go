package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Confing struct {
	Secret        string
	TokenDuration time.Duration
}

type JWTManager struct {
	Confing
}

func NewJWTManager(config *Confing) *JWTManager {
	return &JWTManager{
		*config,
	}
}

func (m *JWTManager) GenerateToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.Secret)

}

// Return claims from signing token
func (m *JWTManager) ParseClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return m.Secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if token.Method != jwt.SigningMethodHS256 {
		return nil, errors.New("invalud signing method")
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (m *JWTManager) GetName(tokenString string) (string, error) {
	claims, err := m.ParseClaims(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}
