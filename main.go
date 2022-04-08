package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	Id          int
	Name        string
	Description string
	Price       int
	Image       string
	CategoryId  int
	Category    Category `gorm:"foreignKey:CategoryId"`
}

type Category struct {
	Id          int
	Name        string
	Description string
}

func main() {
	// Load environment variables from file.
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	// Connect to PlanetScale database using DSN environment variable.
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed to connect to PlanetScale: %v", err)
	}

	// Create a handler whose methods will be used to serve HTTP routes.
	h := &handler{db: db}

	r := mux.NewRouter()
	r.HandleFunc("/seed", h.seedDatabase).Methods(http.MethodGet)
	r.HandleFunc("/products", h.getProducts).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", h.getProduct).Methods(http.MethodGet)
	r.HandleFunc("/categories", h.getCategories).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}", h.getCategory).Methods(http.MethodGet)

	// Start an HTTP API server.
	const addr = ":8080"
	log.Printf("successfully connected to PlanetScale, starting HTTP server on %q", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

type handler struct {
	db *gorm.DB
}

func (h *handler) seedDatabase(w http.ResponseWriter, r *http.Request) {
	h.db.AutoMigrate(&Product{})
	h.db.AutoMigrate(&Category{})

	h.db.Create(&Category{Name: "Phone", Description: "Description 1"})
	h.db.Create(&Category{Name: "Video Game Console", Description: "Description 2"})

	h.db.Create(&Product{Name: "iPhone", Description: "Description 1", Image: "Image 1", Category: Category{Id: 1}})
	h.db.Create(&Product{Name: "Pixel Pro", Description: "Description 2", Image: "Image 2", Category: Category{Id: 1}})
	h.db.Create(&Product{Name: "Playstation", Description: "Description 3", Image: "Image 3", Category: Category{Id: 2}})
	h.db.Create(&Product{Name: "Xbox", Description: "Description 4", Image: "Image 4", Category: Category{Id: 2}})
	h.db.Create(&Product{Name: "Galaxy S", Description: "Description 5", Image: "Image 5", Category: Category{Id: 1}})

	w.Write([]byte("Migrations and Seeding of database complete"))
}

func (h *handler) getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}
	h.db.Preload("Category").Find(&products)
	json.NewEncoder(w).Encode(products)
}

func (h *handler) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product := Product{}
	result := h.db.First(&product, id)

	if result.Error != nil {
		w.Write([]byte("Product not found"))
	} else {
		json.NewEncoder(w).Encode(product)
	}
}

func (h *handler) getCategories(w http.ResponseWriter, r *http.Request) {
	categories := []Category{}
	h.db.Find(&categories)
	json.NewEncoder(w).Encode(categories)
}

func (h *handler) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	category := Category{}
	result := h.db.First(&category, id)

	if result.Error != nil {
		w.Write([]byte("Category not found"))
	} else {
		json.NewEncoder(w).Encode(category)
	}
}
