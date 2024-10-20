package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
)

var token *jwtauth.JWTAuth

var secretKey = []byte("secret-key")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func init() {
	token = jwtauth.New("HS256", []byte("test"), nil)
	_, tokenString, _ := token.Encode(map[string]any{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func main() {
	r := chi.NewRouter()

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var u User
		json.NewDecoder(r.Body).Decode(&u)
		fmt.Printf("The user request value %v", u)

		if u.Username == "Chek" && u.Password == "123456" {
			tokenString, err := CreateToken(u.Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Errorf("No username found")
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, tokenString)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid credentials")
		}
	})

	http.ListenAndServe(":3000", r)
}
