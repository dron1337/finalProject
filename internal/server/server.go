package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dron1337/finalProject/internal/api"
	"github.com/dron1337/finalProject/internal/constants"
	"github.com/joho/godotenv"
)

func Run() {
	 err := godotenv.Load()
	//err := godotenv.Load("C:\\Temp\\Go\\project\\finalProject\\.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	router := api.Init()
	fmt.Printf("Server starts port %s\n", constants.Port)
	log.Fatal(http.ListenAndServe(constants.Port, router))
}
