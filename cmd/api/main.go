package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/konmitnik/hitalent_test/internal/config"
	"github.com/konmitnik/hitalent_test/internal/handlers"
	"github.com/konmitnik/hitalent_test/internal/repository"
)

func main() {
	conn, err := config.NewDBConfig().OpenConnection()
	if err != nil {
		log.Fatalf("db open fail: %v", err)
	}

	repo := repository.NewRepository(conn)

	mux := http.NewServeMux()

	mux.HandleFunc("/departments/", handlers.NewHandler(repo).Handle)

	port := ":8080"
	fmt.Printf("Server started at %s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
