// Command golang-example demonstrates how to connect to PlanetScale from a Go
// application.
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// A Product contains metadata about a product for sale.
type Product struct {
	ID          int
	Name        string
	Description string
	Image       string
	CategoryID  int
	Category    Category `gorm:"foreignKey:CategoryID"`
}

// A Category describes a group of Products.
type Category struct {
	ID          int
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

	// Create an API handler which serves data from PlanetScale.
	handler := NewHandler(db)

	// Start an HTTP API server.
	const addr = ":8080"
	log.Printf("successfully connected to PlanetScale, starting HTTP server on %q", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

// A Handler is an HTTP API server handler.
type Handler struct {
	db *gorm.DB
}

// NewHandler creates an http.Handler which wraps a PlanetScale database
// connection.
func NewHandler(db *gorm.DB) http.Handler {
	h := &Handler{db: db}

	r := mux.NewRouter()
	r.HandleFunc("/seed", h.seedDatabase).Methods(http.MethodGet)
	r.HandleFunc("/products", h.getProducts).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", h.getProduct).Methods(http.MethodGet)
	r.HandleFunc("/categories", h.getCategories).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}", h.getCategory).Methods(http.MethodGet)

	return r
}

// seedDatabase is the HTTP handler for GET /seed.
func (h *Handler) seedDatabase(w http.ResponseWriter, r *http.Request) {
	// Perform initial schema migrations.
	if err := h.db.AutoMigrate(&Product{}); err != nil {
		http.Error(w, "failed to migrate products table", http.StatusInternalServerError)
		return
	}

	if err := h.db.AutoMigrate(&Category{}); err != nil {
		http.Error(w, "failed to migrate categories table", http.StatusInternalServerError)
		return
	}

	// Seed categories and products for those categories.
	h.db.Create(&Category{
		Name:        "Phone",
		Description: "Description 1",
	})
	h.db.Create(&Category{
		Name:        "Video Game Console",
		Description: "Description 2",
	})

	h.db.Create(&Product{
		Name:        "iPhone",
		Description: "Description 1",
		Image:       "Image 1",
		Category:    Category{ID: 1},
	})
	h.db.Create(&Product{
		Name:        "Pixel Pro",
		Description: "Description 2",
		Image:       "Image 2",
		Category:    Category{ID: 1},
	})
	h.db.Create(&Product{
		Name:        "Playstation",
		Description: "Description 3",
		Image:       "Image 3",
		Category:    Category{ID: 2},
	})
	h.db.Create(&Product{
		Name:        "Xbox",
		Description: "Description 4",
		Image:       "Image 4",
		Category:    Category{ID: 2},
	})
	h.db.Create(&Product{
		Name:        "Galaxy S",
		Description: "Description 5",
		Image:       "Image 5",
		Category:    Category{ID: 1},
	})

	io.WriteString(w, "Migrations and Seeding of database complete\n")
}

// getProducts is the HTTP handler for GET /products.
func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product
	result := h.db.Preload("Category").Find(&products)
	if result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(products)
}

// getProduct is the HTTP handler for GET /products/{id}.
func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	result := h.db.First(&product, mux.Vars(r)["id"])
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// getCategories is the HTTP handler for GET /categories.
func (h *Handler) getCategories(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	result := h.db.Find(&categories)
	if result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(categories)
}

// getCategory is the HTTP handler for GET /category/{id}.
func (h *Handler) getCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	result := h.db.First(&category, mux.Vars(r)["id"])
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(category)
}
