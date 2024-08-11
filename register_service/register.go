package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users = make(map[string]User)

func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if _, exists := users[newUser.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	users[newUser.Username] = newUser
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func main() {
	http.HandleFunc("/register", registerUser)
	http.ListenAndServe(":8080", nil)
}