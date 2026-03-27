// kirim json ke server
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Nama string `json:"nama"`
	Umur int    `json:"umur"`
}

type Response struct {
	Message string `json:"message"`
	Data    User   `json:"data"`
}

func main() {
	http.HandleFunc("/user", handler)
	fmt.Println("Server jalan di :8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("harus POST "))
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Write([]byte("Error parsing JSON"))
		return
	}
	if user.Nama == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Nama tidak boleh kosong"))
		return
	}
	if user.Umur <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Umur harus lebih dari 0"))
		return
	}
	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "success",
		Data:    user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
