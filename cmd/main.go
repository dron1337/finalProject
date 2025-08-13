package main

import (
	"log"

	"github.com/dron1337/finalProject/internal/db"
	"github.com/dron1337/finalProject/internal/server"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	server.Run()
}
