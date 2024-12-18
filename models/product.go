package models

import (
	"github.com/guregu/null/v5"
	"time"
)

type (
	Product struct {
		ID        int       `json:"ID,omitempty"`
		SKU       string    `json:"sku"`
		Name      string    `json:"name"`
		Category  string    `json:"category"`
		Price     PriceData `json:"price"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PriceData struct {
		Original           int         `json:"original"`
		Final              int         `json:"final"`
		DiscountPercentage null.String `json:"discount_percentage,omitempty"`
		Currency           string      `json:"currency"`
	}

	ProductsResponse struct {
		Products []Product `json:"products"`
		Meta     Meta      `json:"meta"`
	}

	Products []Product
)
