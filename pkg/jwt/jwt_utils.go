package jwt

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwtlib.RegisteredClaims
}

func VerifyToken(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return &JWTClaims{}, err
	}

	if !token.Valid {
		return &JWTClaims{}, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return &JWTClaims{}, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func GetJwtClaims(r *http.Request) (*JWTClaims, error) {
	claims, ok := r.Context().Value("claims").(*JWTClaims)
	if !ok {
		return &JWTClaims{}, fmt.Errorf("failed to get validated claims")
	}
	return claims, nil
}

func CreateJwtToken(JwtSecretToken string, userID int) (string, error) {
	claims := JWTClaims{
		jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Duration(3600 * time.Second))),
			Subject:   strconv.Itoa(userID),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(JwtSecretToken))
	if err != nil {
		return "", err
	}
	return t, nil
}
