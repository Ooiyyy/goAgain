package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/profile", jwtMiddleware(profileHandler))
	fmt.Println("Server jalan di :8080")
	http.ListenAndServe(":8080", nil)
}

var jwtSecret = []byte("secret123")

func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jsonError(w, http.StatusMethodNotAllowed, "Method harus POST")
		// w.Write([]byte("harus POST "))
	}
	var req LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("Invalid JSON"))
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Username == "" {
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]"Username wajib diisi") (bisa pakai cara ini atau yang atas)
		jsonError(w, http.StatusBadRequest, "Username wajib diisi")
		return
	}
	if req.Password == "" {
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]"Password wajib diisi") (bisa pakai cara ini atau yang atas)
		jsonError(w, http.StatusBadRequest, "Password wajib diisi")
		return
	}
	if req.Username != "admin" || req.Password != "1234" {
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Write([]byte("login gagal"))
		jsonError(w, http.StatusUnauthorized, "Username atau password salah")
		return
	}

	claims := jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, http.StatusUnauthorized, "Token tidak ada")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			jsonError(w, http.StatusUnauthorized, "Token tidak valid")
			return
		}
		next(w, r)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "ini data rahasia",
	})
}
