package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Databaseverbinding opzetten
	var err error
	db, err = sql.Open("mysql", "rik:SQLR1k@tcp(localhost:3306)/Login system DB")
	if err != nil {
		fmt.Println("Fout bij het openen van de database.")
	}
	defer db.Close()
	// Tabel maken als deze nog niet bestaat
	CreateTable()

	// Lees de gebruikersinvoer
	inputUsername, inputPassword := readInput()

	// Haal gebruikersgegevens op
	username, password := getUserData(inputUsername)

	// Controleer de gebruikersnaam
	usernameControl(inputUsername, username)

	// Controleer het ingevoerde wachtwoord
	passwordControl(inputPassword, password)

	fmt.Println(username, password)
}

func CreateTable() {
	// Tabel maken als deze nog niet bestaat
	var err error
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Fout bij het maken van de tabel.")
	}

	fmt.Println("Tabel 'users' is gemaakt of bestaat al.")
}

// Lees de gebruikersinvoer en retourneer deze
func readInput() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Gebruikersnaam: ")
	inputUsername, _ := reader.ReadString('\n')
	fmt.Print("Wachtwoord: ")
	inputPassword, _ := reader.ReadString('\n')

	// Verwijder nieuwe regels van de ingevoerde waarden
	inputUsername = strings.TrimSpace(inputUsername)
	inputPassword = strings.TrimSpace(inputPassword)

	return inputUsername, inputPassword
}

// Query om gebruikersgegevens op te halen en retourneer ze
func getUserData(inputUsername string) (string, string) {
	query := "SELECT * FROM users WHERE username=?"
	row := db.QueryRow(query, inputUsername)

	var username, password string
	err := row.Scan(&username, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Ongeldige gebruikersnaam of wachtwoord.")
		} else {
			fmt.Println("Fout bij het ophalen van de gebruikersgegevens.")
		}
		return "", ""
	}

	return username, password
}

// Controleer het ingevoerde wachtwoord
func passwordControl(inputPassword string, password string) {
	if inputPassword == password {
		fmt.Println("Inloggen gelukt!")
	} else {
		fmt.Println("Ongeldige gebruikersnaam of wachtwoord.")
	}
}

func usernameControl(inputUsername string, username string) {
	if inputUsername == username {
		fmt.Println("Inloggen gelukt!")
	} else {
		fmt.Println("Ongeldige gebruikersnaam of wachtwoord.")
	}
}
