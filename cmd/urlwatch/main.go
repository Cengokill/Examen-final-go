package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Cengokill/Examen-final-go/internal/api"
	"github.com/Cengokill/Examen-final-go/internal/checker"
	"github.com/Cengokill/Examen-final-go/internal/pool"
	"github.com/Cengokill/Examen-final-go/internal/store"
)

func main() {
	logger := api.NewJSONLogger(os.Getenv("LOG_LEVEL"))
	// fmt.Println("LOG_LEVEL =", os.Getenv("LOG_LEVEL")) // test DEBUG vs ERROR

	memStore := store.NewMemoryStore()
	runner := pool.NewRunner(checker.NewHTTPChecker())
	server := api.NewServer(memStore, runner)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	// fmt.Println("URLWatch démarre sur", addr)
	// vérif au démarrage
	// fmt.Println("store + checker HTTP + pool câblés")
	fmt.Println("=== URLWatch ===")
	fmt.Println("Serveur sur http://localhost" + addr)
	fmt.Println("Routes : POST /v1/checks | GET /v1/checks/{id} | GET /healthz")

	if err := http.ListenAndServe(addr, server.Handler(logger)); err != nil {
		// fmt.Println("ListenAndServe crash :", err) // port 8080 déjà pris ?
		log.Fatal(err)
	}
}
