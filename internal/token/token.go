package token

import (
	"errors"
	"fmt"
	"net/http"
	"register-service/internal/config"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func ValidateToken(r *http.Request) error {
	tokenString := getToken(r)

	token, err := jwt.Parse(tokenString, returnSecretKey)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func returnSecretKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signature method %v", token.Header["alg"])
	}

	return []byte(config.Get().Token.Key), nil
}

func ExtractUserId(request *http.Request) (int, error) {
	tokenS, err := jwt.Parse(getToken(request), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid jwt")
		}

		return []byte(config.Get().Token.Key), nil
	})
	if err != nil {
		return 0, err
	}

	return tokenS.Claims.(jwt.MapClaims)["userId"].(int), nil
}
