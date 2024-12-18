package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/tonymj76/mytheresa-test/ent"

	"log"

	_ "github.com/lib/pq"
)

func NewDB(link string) (*ent.Client, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"),
		os.Getenv("DB_POST"), os.Getenv("DB_SSLMODE"))
	if link == "" {
		link = dsn
	}
	client, err := ent.Open("postgres", link)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		defer client.Close()
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client, nil
}
