package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBname string `json:"dbname"`
	DBuser string `json:"dbuser"`
	DBpass string `json:"dbpass"`
	DBhost string `json:"dbhost"`
	DBport string `json:"dbport"`
}

var (
	db        *sql.DB
	templates *template.Template
)

func main() {
	// Read the config file
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	dsn := fmt.Sprintf("%s:%s@/%s", config.DBuser, config.DBpass, config.DBname)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the necessary database table if it doesn't exist
	if err := CreateTable(); err != nil {
		log.Fatal(err)
	}

	// Load HTML templates
	templates = template.Must(template.ParseGlob("templates/*.html"))

	// Register HTTP handlers
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)

	// Start the web server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Println("error:", err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if err := checkCredentials(username, password); err == nil {
			fmt.Fprintln(w, "Login successful!")
			return
		}

		fmt.Fprintln(w, "Invalid username or password.")
		return
	}

	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		log.Println("error:", err)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if err := addUser(username, password); err == nil {
			fmt.Fprintln(w, "Registration successful!")
			return
		}

		fmt.Fprintln(w, "Username is already in use.")
		return
	}

	if err := templates.ExecuteTemplate(w, "register.html", nil); err != nil {
		log.Println("error:", err)
	}
}

func CreateTable() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

// Voeg een gebruiker toe aan de database
func addUser(username, password string) error {
	// Controleer of de gebruiker al bestaat in de database
	query := "SELECT COUNT(*) FROM users WHERE username=?"
	row := db.QueryRow(query, username)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("gebruikersnaam is al in gebruik")
	}

	// Voeg de gebruiker toe aan de database
	insertQuery := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = db.Exec(insertQuery, username, password)
	if err != nil {
		return err
	}

	return nil
}

// Controleer de gebruikersnaam en het wachtwoord
func checkCredentials(username, password string) error {
	query := "SELECT COUNT(*) FROM users WHERE username=? AND password=?"
	row := db.QueryRow(query, username, password)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("ongeldige gebruikersnaam of wachtwoord")
	}

	return nil
}
