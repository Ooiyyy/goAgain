package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UpdateUserReq struct {
	Password string `json:"password"`
}

func main() {
	connectDB()
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/profile", jwtMiddleware(profileHandler))
	http.HandleFunc("/update", jwtMiddleware(updateUserHandler))
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

	var username string
	var password string

	query := "SELECT username, password FROM users where username = ?"
	err = db.QueryRow(query, req.Username).Scan(&username, &password)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "User tidak ditemukan")
	}
	// if req.Password != password {
	// 	jsonError(w, http.StatusUnauthorized, "Password salah")
	// }
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password))
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "password salah")
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
	// if req.Username != "admin" || req.Password != "1234" {
	// 	// w.WriteHeader(http.StatusUnauthorized)
	// 	// w.Write([]byte("login gagal"))
	// 	jsonError(w, http.StatusUnauthorized, "Username atau password salah")
	// 	return
	// }

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
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "ini data rahasia",
		"username": username,
	})
}

var db *sql.DB

func connectDB() {
	dsn := "root:Gridy20@tcp(127.0.0.1:3306)/testdb"

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = database.Ping()
	if err != nil {
		panic(err)
	}

	db = database
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if r.Method != "POST" {
		jsonError(w, http.StatusMethodNotAllowed, "Method harus POST")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if req.Username == "" {
		jsonError(w, http.StatusBadRequest, "Username wajib diisi")
		return
	}
	if req.Password == "" {
		jsonError(w, http.StatusBadRequest, "Password wajib diisi")
		return
	}
	if len(req.Password) < 6 {
		jsonError(w, http.StatusBadRequest, "password minimal 6 karakter")
		return
	}
	var existing string
	err = db.QueryRow("SELECT username FROM users WHERE username = ?", req.Username).Scan(&existing)

	if err == nil {
		jsonError(w, http.StatusBadRequest, "username sudah digunakan")
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal hash Password")
	}

	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = db.Exec(query, req.Username, string(hashedPass))
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal insert user")
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user berhasil dibuat",
	})
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		jsonError(w, http.StatusMethodNotAllowed, "Method harus PUT")
		return
	}
	username := r.Context().Value("username").(string)

	var req UpdateUserReq
	err := json.NewDecoder(r.Body).Decode((&req))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Password == "" {
		jsonError(w, http.StatusBadRequest, "password wajib diisi")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal hash password")
		return
	}
	query := "UPDATE users SET password = ? WHERE username = ?"
	_, err = db.Exec(query, string(hashedPassword), username)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "gagal update user")
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "password berhasil diupdate",
	})
}
