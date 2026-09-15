package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"recipea.com/m/config"
)

type JWTProvider struct {
	Secret []byte
}

func NewJWTProvider(secret string) *JWTProvider {
	if secret == "" {
		secret = "default-secret-key"
	}
	return &JWTProvider{Secret: []byte(secret)}
}

func (j *JWTProvider) CreateToken(subject uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}

func (j *JWTProvider) ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return j.Secret, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"]
		if !ok {
			return 0, errors.New("missing subject")
		}
		idf, ok := sub.(float64)
		if !ok {
			return 0, errors.New("invalid subject type")
		}
		return uint(idf), nil
	}
	return 0, errors.New("invalid token")
}

func DefaultSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return config.Load().JWTSecret
}
