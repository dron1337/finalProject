package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dron1337/finalProject/internal/api"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()
	// err := godotenv.Load("C:\\Temp\\Go\\project\\finalProject\\.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	router := api.Init()
	fmt.Printf("Server starts port %s\n", os.Getenv("TODO_PORT"))
	log.Fatal(http.ListenAndServe(os.Getenv("TODO_PORT"), router))
}
