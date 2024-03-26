package jwtgen

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v4"

	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

type JWTManager struct {
	TokenName string
	secretKey string
	tokenExp  time.Duration
}

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

// NewJWTManager returns a new instance of JWTManager.
func NewJWTManager(tokenName string, secretKey string, hours int) *JWTManager {
	j := &JWTManager{
		TokenName: tokenName,
		secretKey: secretKey,
		tokenExp:  time.Duration(hours * int(time.Hour)),
	}
	return j
}

// BuildJWTString creates JWT token with userID.
func (j *JWTManager) BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID returns userID from JWT token.
func (j *JWTManager) GetUserID(tokenString string) (string, error) {
	jwtClaims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperr.NewValueError(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]), apperr.Caller(), errors.New("unexpected signing method"))
			}
			return []byte(j.secretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		slog.Warn("token is not valid", slog.String("error", err.Error()))
		return "", apperr.NewValueError("token is not valid", apperr.Caller(), errors.New("token is not valid"))
	}

	return jwtClaims.UserID, nil
}
