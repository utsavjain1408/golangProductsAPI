package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

//COLLECTION is the collection in the mongoDB database which will store our products documents
const COLLECTION = "Products"

//Get all products
func getProducts(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()
		db := session.DB("ProductCatelogue").C(COLLECTION)
		w.Header().Set("Content-Type", "application/json")
		var products []Product
		err := db.Find(bson.M{}).All(&products)
		if err != nil {
			log.Fatal(err)
			json.NewEncoder(w).Encode(&Product{})

		}
		json.NewEncoder(w).Encode(products)
	}
}

//Get a particular product by ID
func getProduct(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer s.Close()
		db := session.DB("ProductCatelogue").C(COLLECTION)
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r) //Get params
		//Loop through the products and find with ProductID
		//for _, item := range p {
		//	if item.ProductID == params["id"] {
		//		json.NewEncoder(w).Encode(item)
		//		return
		//	}
		//}
		//if not found then
		var product Product
		err := db.FindId(bson.ObjectIdHex(params["id"])).One(&product)
		if err != nil {
			log.Fatal(err)
			json.NewEncoder(w).Encode(&Product{})

		}
		json.NewEncoder(w).Encode(product)
	}
}

//Inserting a Product in the database
func addProduct(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()
		db := session.DB("ProductCatelogue").C(COLLECTION)
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()
		var pro Product
		if err := json.NewDecoder(r.Body).Decode(&pro); err != nil {
			_ = json.NewDecoder(r.Body).Decode(&pro)
			return
		}
		pro.ProductID = strconv.Itoa(rand.Intn(1000)) //Creates mock product ID
		err := db.Insert(pro)
		if err != nil {
			log.Fatal(err)
			json.NewEncoder(w).Encode(&Product{})
		}
		json.NewEncoder(w).Encode(pro)
	}
}

//Changing some property of a particular product
func updateProduct(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		/*	for index, item := range p {
			if item.ProductID == params["id"] {
				p = append(p[:index], p[index+1:]...)
				var pro Product
				_ = json.NewDecoder(r.Body).Decode(&pro)
				pro.ProductID = params["id"] //Creates mock product ID
				p = append(p, pro)
				json.NewEncoder(w).Encode(pro)
				return
			}
		}*/
		var pro Product
		_ = json.NewDecoder(r.Body).Decode(&pro)
		json.NewEncoder(w).Encode(pro)
		db := session.DB("ProductCatelogue").C(COLLECTION)
		err := db.Update(bson.M{"_id": params["id"]}, &pro)
		if err != nil {
			log.Println("Error!!!!")
		}
	}
}

func deleteProduct(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		session := s.Copy()
		defer session.Close()
		params := mux.Vars(r)
		//for index, item := range p {
		//	if item.ProductID == params["id"] {
		//		p = append(p[:index], p[index+1:]...)
		//		break
		//	}
		//}
		//err := db.Remove(&pro)
		db := session.DB("ProductCatelogue").C(COLLECTION)
		err := db.Remove(params["id"])
		if err != nil {
			log.Println("Error!!!!" + params["id"])
		}
		json.NewEncoder(w).Encode(params["id"])
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(`This is home page, Welcome`)
}

//Person This was for sample mongodbTest
type Person struct {
	Name  string
	Phone string
}

func main() {

	// Init Router
	r := mux.NewRouter()
	session, err := mgo.Dial("35.194.50.254:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

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
	r.HandleFunc("/api/products", getProducts(session)).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct(session)).Methods("GET")
	r.HandleFunc("/api/products", addProduct(session)).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct(session)).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct(session)).Methods("DELETE")
	r.HandleFunc("/", home).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))

}
