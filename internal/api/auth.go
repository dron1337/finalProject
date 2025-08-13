package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/securecookie"
)

const (
	tokenExpiration = 8 * time.Hour
	cookieMaxAge    = 480
)

var (
	hashKey  = securecookie.GenerateRandomKey(64)
	blockKey = securecookie.GenerateRandomKey(32)
	s        = securecookie.New(hashKey, blockKey)
)

type Claims struct {
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("TODO_PASSWORD")
		if len(secret) > 0 {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Error(w, "Authorization cookie missing", http.StatusUnauthorized)
				return
			}
			tokenString := cookie.Value
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil || token == nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func generateToken() (string, error) {
	secret := []byte(os.Getenv("TODO_PASSWORD"))
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "scheduler",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func SetToken(passUser string, w http.ResponseWriter) error {
	passSystem := os.Getenv("TODO_PASSWORD")
	if passSystem != passUser {
		return fmt.Errorf("Incorrect password")
	}
	token, err := generateToken()
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieMaxAge,
	}
	http.SetCookie(w, cookie)
	return nil
}
