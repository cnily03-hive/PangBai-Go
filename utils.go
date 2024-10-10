package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Stringer struct{}

func (s Stringer) String() string {
	return "[struct]"
}

func RandString(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}
	return string(b)
}

func genJwt(o Token) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": o.Name,
		"exp":  time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.JwtKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func validateJwt(tokenStr string) (Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtKey), nil
	})
	if err != nil {
		return Token{}, err
	}
	var user *Token = nil
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		un := claims["user"].(string)
		user = &Token{Name: un}
		return *user, err
	} else {
		// Invalid token
		return Token{}, err
	}
}

// This function disables directory listing for http.FileServer
func noDirList(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") || r.URL.Path == "" {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
