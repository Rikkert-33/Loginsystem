package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	OpenDB()
	readInput()
	getUserData(inputUsername)
	passwordControl(inputPassword, password)

}

// Databaseverbinding opzetten
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "rik:SQLR1k@tcp(localhost:3306)/Login system DB")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

// Lees de gebruikersinvoer
func readInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Gebruikersnaam: ")
	inputUsername, _ := reader.ReadString('\n')
	fmt.Print("Wachtwoord: ")
	inputPassword, _ := reader.ReadString('\n')

	// Verwijder nieuwe regels van de ingevoerde waarden
	inputUsername = strings.TrimSpace(inputUsername)
	inputPassword = strings.TrimSpace(inputPassword)
}

// Query om gebruikersgegevens op te halen
func getUserData(inputUsername string) {
	query := "SELECT * FROM users WHERE username=?"
	row := db.QueryRow(query, inputUsername)

	var username, password string
	err = row.Scan(&username, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Ongeldige gebruikersnaam of wachtwoord.")
		} else {
			panic(err)
		}
		return
	}
}

// Controleer het ingevoerde wachtwoord
func passwordControl(inputPassword string, password string) {
	if inputPassword == password {
		fmt.Println("Inloggen gelukt!")
	} else {
		fmt.Println("Ongeldige gebruikersnaam of wachtwoord.")
	}
}
