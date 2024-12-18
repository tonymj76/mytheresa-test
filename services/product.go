package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v5"
	"github.com/tonymj76/mytheresa-test/ent"
	"github.com/tonymj76/mytheresa-test/ent/category"
	"github.com/tonymj76/mytheresa-test/ent/product"
	"github.com/tonymj76/mytheresa-test/models"
)

// discountRecord is a map that store category name or sku with the percentage discount to be associated with it.
// for more scalability in production discountRecord will have its on table.
var discountRecord = map[string]float64{
	"boots":  0.30, // 30%
	"000003": 0.15, // 15%
}

const CURRENCY = "EUR"

func applyDiscount(epd *ent.Product) models.Product {
	var discount float64
	var pd models.Product
	if value, ok := discountRecord[epd.Edges.Category.Name]; ok {
		discount = value
	}

	if value, ok := discountRecord[epd.Sku]; ok {
		discount = max(discount, value)
	}
	if discount > 0 {
		pd.Price.DiscountPercentage = null.StringFrom(fmt.Sprintf("%v%%", discount*100))
		pd.Price.Final = int(float64(epd.Price) * (1 - discount))
	} else {
		pd.Price.Final = epd.Price
	}
	return pd
}

func applyResponseFields(epd *ent.Product) models.Product {
	pd := applyDiscount(epd)
	pd.Price.Original = epd.Price
	pd.Price.Currency = CURRENCY
	pd.ID = epd.ID
	pd.SKU = epd.Sku
	pd.Name = epd.Name
	pd.CreatedAt = epd.CreatedAt
	pd.UpdatedAt = epd.UpdatedAt
	pd.Category = epd.Edges.Category.Name
	return pd
}

// FilterProduct help to filter product base on category or price less than the value provide
func (rs *RestService) FilterProduct(c *gin.Context, cate string, priceLessThan, page, limit int) (*models.ProductsResponse, error) {
	var dbProducts []*ent.Product
	var dbError error
	var products models.Products

	// Calculate offset
	offset := (page - 1) * limit

	// Get total count of products
	total, err := rs.DB.Product.Query().Count(c)
	if err != nil {
		return nil, fmt.Errorf("failed counting products: %w", err)
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	// Query products with pagination
	switch {
	case cate != "":
		dbProducts, dbError = rs.DB.Product.Query().
			WithCategory().
			Where(product.HasCategoryWith(category.Name(cate))).
			Limit(limit).
			Offset(offset).
			All(c)
	case priceLessThan > 0:
		dbProducts, dbError = rs.DB.Product.Query().WithCategory().
			WithCategory().
			Where(product.PriceLTE(priceLessThan)).
			Limit(limit).
			Offset(offset).
			All(c)
	default:
		dbProducts, dbError = rs.DB.Product.Query().WithCategory().
			WithCategory().
			Limit(limit).
			Offset(offset).
			All(c)
	}

	if dbError != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", dbError)
	}

	for _, dbProduct := range dbProducts {
		products = append(products, applyResponseFields(dbProduct))
	}

	// Build response
	response := &models.ProductsResponse{
		Products: products,
		Meta: models.Meta{
			TotalRecords: total,
			Page:         page,
			TotalPages:   totalPages,
			Limit:        limit,
		},
	}

	return response, nil
}
