package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mamenzul/go-api/cmd/api"
	"github.com/mamenzul/go-api/configs"
	"github.com/mamenzul/go-api/db"
)

func main() {
	db, err := db.CreateDb(configs.Envs.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
