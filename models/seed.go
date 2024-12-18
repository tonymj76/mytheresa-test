package models

type SeedData struct {
	Categories []CategorySeed `json:"categories"`
	Products   []ProductSeed  `json:"products"`
}

type CategorySeed struct {
	Name string `json:"name"`
}

type ProductSeed struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}
