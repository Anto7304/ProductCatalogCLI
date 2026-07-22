package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" //_ mean import for side effects o
)

func main() {

	fmt.Println("Starting product catalog cli.....")
	fmt.Println("Attempting to connect to postgreSQL....")

	connstr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	//ATTEMPT TO OPEN THE CONNECTIO

	db, err := sql.Open("postgres", connstr)

	if err != nil {
		log.Fatal("Failed to open database connection", err)

	}

	defer db.Close() //ensure connection closes when program ends

	//test connection

	err = db.Ping()

	if err != nil {
		log.Fatal("failed to ping database")
	}

	fmt.Println("Successful connected to postgreSQL!")
	fmt.Printf("Connected to database: %s\n", os.Getenv("DB_NAME"))

}
