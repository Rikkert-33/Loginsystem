package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Databaseverbinding opzetten
	var err error
	db, err = sql.Open("mysql", "rik:SQLR1k@tcp(localhost:3306)/mijndb2")
	if err != nil {
		fmt.Println("Fout bij het openen van de database.", err)
		return
	}
	defer db.Close()

	// Tabel maken als deze nog niet bestaat
	CreateTable()

	var choice int
	fmt.Println("Welkom bij de inlogapplicatie! Wilt u inloggen of registreren?")
	fmt.Println("Type 1 voor inloggen of 2 voor registreren:")
	fmt.Scanln(&choice)

	if choice == 1 {
		// Lees de gebruikersinvoer
		inputUsername, inputPassword := readInput()
		// Controleer de gebruikersnaam en het wachtwoord
		err := checkCredentials(inputUsername, inputPassword)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Inloggen gelukt!")
	} else if choice == 2 {
		// Lees de gebruikersinvoer
		inputUsername, inputPassword := readInput()

		// Voeg de gebruiker toe aan de database
		err := addUser(inputUsername, inputPassword)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Gebruiker is toegevoegd aan de database.")
	} else {
		fmt.Println("Ongeldige keuze.")
	}

}

func CreateTable() error {
	// Tabel maken als deze nog niet bestaat
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

// Query om gebruikersgegevens op te halen en retourneer ze
// func getUserData(inputUsername string) (inputPassword, string) (string, string) {
// 	query := "SELECT username, password FROM users WHERE username=?"
// 	row := db.QueryRow(query, inputUsername, inputPassword)

// 	var username, password string
// 	err := row.Scan(&username, &password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("Gebruiker niet gevonden.")
// 		} else {
// 			fmt.Println("Fout bij het ophalen van de gebruikersgegevens.", err)
// 		}
// 		return "", ""
// 	}

// 	fmt.Println("Inloggen gelukt!")
// 	return username, password

// }

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
