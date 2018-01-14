package main

type Product struct {
	ProductID   string `json:"id"`
	Name        string `json:"name"`
	SKU         string `json:"SKU"`
	Description string `json:"description"`
	Store       *Store `json:"store"`
}

//Store specific details of a product
type Store struct {
	StoreName string  `json:"storeName"`
	Price     float32 `json:"price"`
	Stock     int     `json:"stock"`
}
