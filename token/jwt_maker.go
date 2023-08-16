package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTmaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: %v", minSecretKeySize)
	}
	return &JWTmaker{secretKey: secretKey}, nil
}

// Creates token for a specific username
func (maker *JWTmaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":         payload.ID,
		"username":   payload.Username,
		"issued_at":  payload.IssuedAt.Unix(),
		"expired_at": payload.ExpiresAt.Unix(),
	}

	payload = nil

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// Check if the token is valid
func (maker *JWTmaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims := jwtToken.Claims.(jwt.MapClaims)
	if jwtToken.Valid {
		uid := uuid.MustParse(claims["id"].(string))

		payload := &Payload{
			ID:        uid,
			Username:  claims["username"].(string),
			IssuedAt:  time.Unix(int64(claims["issued_at"].(float64)), 0),
			ExpiresAt: time.Unix(int64(claims["expired_at"].(float64)), 0),
		}

		if err = payload.Valid(); err != nil {
			return nil, err
		}

		return payload, nil
	}
	return nil, ErrInvalidToken
}
