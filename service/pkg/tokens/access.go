package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mewov/authorization-rest/models"
)

type (
	Claims struct {
		UserId int64  `json:"user_id"`
		Login  string `json:"login"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}
)

func GenerateAccess(config *models.LocalConfig, userId int64, login string, email string, role string) (string, error) {
	claims := Claims{
		UserId: userId, Login: login,
		Email: email, Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SERVER_PASSWORD))
}

func CheckAccess(config *models.LocalConfig, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SERVER_PASSWORD), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
