package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var (
	db  *gorm.DB
	err error
)

// Product struct
type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

// Result struct
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/simple_mux?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connection to database failed")
	} else {
		log.Println("Connection to database success")
	}

	db.AutoMigrate(&Product{})
	handleRequest()
}

func handleRequest() {
	log.Println("Application running on PORT 3000")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/product", createProduct).Methods("POST")
	router.HandleFunc("/products", getProducts).Methods("GET")
	router.HandleFunc("/product/{id}", getProduct).Methods("GET")
	router.HandleFunc("/product/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Create product success",
	}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)

	res := Result{
		Code:    200,
		Data:    products,
		Message: "Ger product success",
	}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Ger product success",
	}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var productUpdates Product
	json.Unmarshal(payloads, &productUpdates)

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdates)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Update product success",
	}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)
	db.Delete(&product)

	res := Result{
		Code:    200,
		Message: "Delete product success",
	}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
