package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hymiside/fitness-api/pkg/handler"
	"github.com/Hymiside/fitness-api/pkg/repository"
	"github.com/Hymiside/fitness-api/pkg/service"
	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Panicf("error to load .env file: %v", err)
	}

	psqlInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", 
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DATABASE"),
	)
	db, err := dbr.Open("postgres", psqlInfo, nil)
	if err != nil {
		log.Fatalf("failed to connection postgres: %v", err)
	}

	log.Println("connection postgres test...")
	if err = db.Ping(); err != nil {
		log.Fatalf("connection test error: %v", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)


	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-quit:
			cancel()
		case <-ctx.Done():
			db.Close()
			return
		}
	}()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
		Handler: handlers.InitRoutes(),
	}

	go func() {
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()

	log.Printf("server started on http://%s/", httpServer.Addr)
	if err = httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return
		}
		log.Fatalf("failed to start server: %v", err)
	}
}