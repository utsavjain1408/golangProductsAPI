package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Product model for a Product ...
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

var p []Product

//Get all products
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)

}

//Get a particular product by ID
func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //Get params
	//Loop through the products and find with ProductID
	for _, item := range p {
		if item.ProductID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	//if not found then
	json.NewEncoder(w).Encode(&Product{})

}

//Inserting a Product in the database
func addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var pro Product
	_ = json.NewDecoder(r.Body).Decode(&pro)
	pro.ProductID = strconv.Itoa(rand.Intn(1000)) //Creates mock product ID
	p = append(p, pro)
	json.NewEncoder(w).Encode(pro)

}

//Changing some property of a particular product
func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range p {
		if item.ProductID == params["id"] {
			p = append(p[:index], p[index+1:]...)
			var pro Product
			_ = json.NewDecoder(r.Body).Decode(&pro)
			pro.ProductID = params["id"] //Creates mock product ID
			p = append(p, pro)
			json.NewEncoder(w).Encode(pro)
			return
		}
	}
	json.NewEncoder(w).Encode(p)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range p {
		if item.ProductID == params["id"] {
			p = append(p[:index], p[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(p)
}

func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(`This is home page, Welcome`)
}

// Sample Data

func main() {
	// Init Router
	r := mux.NewRouter()

	// Sample Data
	p = append(p, Product{
		ProductID:   "3",
		Name:        "Coke Cola",
		SKU:         "dozen",
		Description: "This a some product",
		Store: &Store{
			StoreName: "Target",
			Price:     5,
			Stock:     11}})
	p = append(p, Product{
		ProductID:   "4",
		Name:        "Srite",
		SKU:         "dozen",
		Description: "This a some other product",
		Store: &Store{
			StoreName: "Walmart",
			Price:     4.5,
			Stock:     10}})
	//Route Handlers
	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/products", addProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	r.HandleFunc("/", home).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))

}
