package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBname string `json:"dbname"`
	DBuser string `json:"dbuser"`
	DBpass string `json:"dbpass"`
	DBhost string `json:"dbhost"`
	DBport string `json:"dbport"`
}

var db *sql.DB

func main() {
	//Read the config file
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Databaseverbinding opzetten
	dsn := fmt.Sprintf("%s:%s@/%s", config.DBuser, config.DBpass, config.DBname)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
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

	// create the gin router
	router := gin.Default()

	// define the route
	router.GET("/", indexHandler)
	router.POST("/login", loginHandler)
	router.POST("/register", registerHandler)

	// run the server
	router.Run("localhost:8080")

}

// LoadConfig loads the configuration from a JSON file
func loadConfig(filename string) (Config, error) {
	config := Config{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check the credentials
	err := checkCredentials(username, password)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "success.html", nil)
}

func registerHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Add the user to the database
	err := addUser(username, password)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "success.html", nil)
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
