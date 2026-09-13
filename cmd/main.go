package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	srv := &http.Server{
		Addr:              listenAddr(),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("LearnGO listen %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	case sig := <-sigCh:
		log.Printf("stop signal %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}

func listenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("HOST")
	if host == "" && os.Getenv("RENDER") != "" {
		host = "0.0.0.0"
	}
	if host == "" {
		return ":" + port
	}
	return host + ":" + port
}
