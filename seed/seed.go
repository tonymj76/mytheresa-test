package seed

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tonymj76/mytheresa-test/ent"
	"github.com/tonymj76/mytheresa-test/ent/category"
	"github.com/tonymj76/mytheresa-test/ent/product"
	"github.com/tonymj76/mytheresa-test/models"
	"os"
)

// SeedDatabase seeds the database with initial data
func SeedDatabase(client *ent.Client, filepath string) error {
	ctx := context.Background()

	// Load JSON file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %v", err)
	}

	var seedData models.SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Seed categories
	categoryMap := make(map[string]*ent.Category)
	for _, cat := range seedData.Categories {
		existingCategory, err := client.Category.
			Query().
			Where(category.NameEQ(cat.Name)).
			Only(ctx)

		if err != nil { // Category doesn't exist; create it
			log.Printf("Creating category: %s", cat.Name)
			newCategory, err := client.Category.Create().SetName(cat.Name).SetDescription("best product ever").Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create category %s: %v", cat.Name, err)
			}
			categoryMap[cat.Name] = newCategory
		} else {
			log.Printf("Category already exists: %s", cat.Name)
			categoryMap[cat.Name] = existingCategory
		}
	}

	// Seed products
	for _, prod := range seedData.Products {
		cate, exists := categoryMap[prod.Category]
		if !exists {
			return fmt.Errorf("category %s not found for product %s", prod.Category, prod.SKU)
		}

		// Check if product already exists by SKU
		_, err := client.Product.
			Query().
			Where(product.Sku(prod.SKU)).
			Only(ctx)
		if err == nil { // Product already exists
			log.Printf("Product already exists: %s", prod.SKU)
			continue
		}

		log.Printf("Creating product: %s", prod.SKU)
		_, err = client.Product.
			Create().
			SetSku(prod.SKU).
			SetName(prod.Name).
			SetPrice(prod.Price).
			SetCategory(cate).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create product %s: %v", prod.SKU, err)
		}
	}

	return nil
}
