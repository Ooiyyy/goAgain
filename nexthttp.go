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

func main() {
	http.HandleFunc("/user", handler)
	fmt.Println("server berjalan di port :8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := User{
		Nama: "Helmi",
		Umur: 23,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
