package main

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tonymj76/mytheresa-test/config"
	"github.com/tonymj76/mytheresa-test/handlers"
	"github.com/tonymj76/mytheresa-test/routes"
	"github.com/tonymj76/mytheresa-test/seed"
	"github.com/tonymj76/mytheresa-test/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// init gets called before the main function
func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found create")
		os.Exit(1)
	}
}

func main() {
	service, err := services.NewRestService(services.WithDBSetup())
	if err != nil {
		log.Fatalf("error setting up new rest server. Err: %v", err)
	}

	if err := seed.SeedDatabase(service.DB, "seed-product-and-category.json"); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	handler := handlers.NewRegisteredHandler(service)
	route := routes.SetRouter(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.GetEnv("PORT", "9090")),
		Handler: route,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
