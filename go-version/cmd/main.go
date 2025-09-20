package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-version/internal/api"
	"go-version/internal/dal"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Error loading .env file. Continuing with system environment variables.")
	}

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	db, err := dal.NewDatabaseConn(ctx, "./database.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := chi.NewRouter()
	apiService := api.NewService(db)
	apiService.RegisterRoutes(router)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
			done <- true
		}
	}()

	go func() {
		sig := <-sigs
		server.Shutdown(ctx)
		fmt.Printf("Received signal %v, shutting down...\n", sig)
		done <- true
	}()

	fmt.Printf("Server is running on %s\n", fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
	<-done
	fmt.Println("Program exited")
}
