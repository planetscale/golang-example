package main

import (
	"encoding/json"
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

var DB *gorm.DB

func init() {
	// Load environment variables
	godotenv.Load()

	// Connect to PlanetScale database
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db
}

func main() {

	r := mux.NewRouter()
	// Run migrations and load sample data
	r.HandleFunc("/seed", seedDatabase).Methods("GET")
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/categories", getCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", getCategory).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func seedDatabase(w http.ResponseWriter, r *http.Request) {
	DB.AutoMigrate(&Product{})
	DB.AutoMigrate(&Category{})

	DB.Create(&Category{Name: "Phone", Description: "Description 1"})
	DB.Create(&Category{Name: "Video Game Console", Description: "Description 2"})

	DB.Create(&Product{Name: "iPhone", Description: "Description 1", Image: "Image 1", Category: Category{Id: 1}})
	DB.Create(&Product{Name: "Pixel Pro", Description: "Description 2", Image: "Image 2", Category: Category{Id: 1}})
	DB.Create(&Product{Name: "Playstation", Description: "Description 3", Image: "Image 3", Category: Category{Id: 2}})
	DB.Create(&Product{Name: "Xbox", Description: "Description 4", Image: "Image 4", Category: Category{Id: 2}})
	DB.Create(&Product{Name: "Galaxy S", Description: "Description 5", Image: "Image 5", Category: Category{Id: 1}})

	w.Write([]byte("Migrations and Seeding of database complete"))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}
	DB.Preload("Category").Find(&products)
	json.NewEncoder(w).Encode(products)
}
func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product := Product{}
	result := DB.First(&product, id)

	if result.Error != nil {
		w.Write([]byte("Product not found"))
	} else {
		json.NewEncoder(w).Encode(product)
	}

}
func getCategories(w http.ResponseWriter, r *http.Request) {
	categories := []Category{}
	DB.Find(&categories)
	json.NewEncoder(w).Encode(categories)
}
func getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	category := Category{}
	result := DB.First(&category, id)

	if result.Error != nil {
		w.Write([]byte("Category not found"))
	} else {
		json.NewEncoder(w).Encode(category)
	}
}
