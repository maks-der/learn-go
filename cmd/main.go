package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"learngo/internal/content"
	"learngo/internal/web"
)

func main() {
	root, err := content.FindRoot()
	if err != nil {
		log.Fatal(err)
	}
	store, err := content.NewStore(root)
	if err != nil {
		log.Fatal(err)
	}
	handler, err := web.New(store, root)
	if err != nil {
		log.Fatal(err)
	}

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("LearnGO http://localhost%s", addr)
	log.Fatal(srv.ListenAndServe())
}
