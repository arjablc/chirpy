package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	},
	)
	secretBytes := []byte(tokenSecret)
	return token.SignedString(secretBytes)

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(t *jwt.Token) (any, error) {
			byteSec := []byte(tokenSecret)
			return byteSec, nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	uid, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, err
	}
	return uid, nil

}
