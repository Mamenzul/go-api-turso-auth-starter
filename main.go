package main

import (
	"log"
	"net/http"

	"github.com/mamenzul/go-api/db"
	"github.com/mamenzul/go-api/handlers/auth"
	"github.com/mamenzul/go-api/handlers/products"
	"github.com/mamenzul/go-api/middleware"
)

func main() {
	db := db.CreateDb()
	router := http.NewServeMux()
	router.Handle("/products", products.Router(db))
	router.Handle("/", auth.Router(db))
	chain := middleware.MiddlewareChain(middleware.Logger, middleware.JSONMiddleware)

	//create server
	server := &http.Server{
		Addr:    ":8080",
		Handler: chain(router),
	}

	log.Println("Server started at :8080")
	server.ListenAndServe()
	defer db.Close()
}
