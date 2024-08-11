package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *sql.DB

func initDB() {
	var err error
	// Ustal połączenie do bazy danych
	connStr := "user=yourusername dbname=pizzeria sslmode=disable" // Zaktualizuj z odpowiednimi danymi
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Sprawdzenie połączenia z bazą danych
	if err = db.Ping(); err != nil {
		panic(err)
	}
}

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

	// Sprawdzenie, czy użytkownik już istnieje
	var existingUser User
	err = db.QueryRow("SELECT username, email FROM users WHERE username = $1", newUser.Username).Scan(&existingUser.Username, &existingUser.Email)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existingUser.Username != "" {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Haszowanie hasła
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	// Dodanie nowego użytkownika do bazy danych
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Wysłanie odpowiedzi
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Pobranie użytkownika z bazy danych
	var user User
	err = db.QueryRow("SELECT username, email, password FROM users WHERE email = $1", loginReq.Email).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Sprawdzenie hasła
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Wysłanie pozytywnej odpowiedzi
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Login successful")
}

func main() {
	initDB() // Inicjalizacja bazy danych
	defer db.Close() // Zamknij połączenie z bazą danych po zakończeniu

	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8080", nil)
}