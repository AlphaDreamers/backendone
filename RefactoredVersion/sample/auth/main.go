package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Routes
	r.Route("/products", func(r chi.Router) {
		r.Get("/", listProducts)                             // GET /products?category=electronics&page=1
		r.Post("/", createProduct)                           // POST /products
		r.Get("/search", searchProducts)                     // GET /products/search?q=phone&min_price=100&max_price=1000
		r.Get("/category/{category}", getProductsByCategory) // GET /products/category/electronics?sort=price

		r.Route("/{productID}", func(r chi.Router) {
			r.Get("/", getProduct)                // GET /products/123
			r.Put("/", updateProduct)             // PUT /products/123
			r.Patch("/", patchProduct)            // PATCH /products/123
			r.Delete("/", deleteProduct)          // DELETE /products/123
			r.Get("/similar", getSimilarProducts) // GET /products/123/similar?limit=3
		})
	})

	port := ":3002"
	log.Printf("Simple product service running on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

// Handlers with string responses
func listProducts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	page := r.URL.Query().Get("page")

	response := fmt.Sprintf("Listing products - Category: %s, Page: %s", category, page)
	w.Write([]byte(response))
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	category := r.URL.Query().Get("category")
	price := r.URL.Query().Get("price")

	response := fmt.Sprintf("Created product - Name: %s, Category: %s, Price: %s", name, category, price)
	w.Write([]byte(response))
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	response := fmt.Sprintf("Product details for ID: %s", productID)
	w.Write([]byte(response))
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	name := r.URL.Query().Get("name")
	price := r.URL.Query().Get("price")

	response := fmt.Sprintf("Updated product ID: %s - Name: %s, Price: %s", productID, name, price)
	w.Write([]byte(response))
}

func patchProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	field := r.URL.Query().Get("field")
	value := r.URL.Query().Get("value")

	response := fmt.Sprintf("Patched product ID: %s - Field: %s, Value: %s", productID, field, value)
	w.Write([]byte(response))
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	response := fmt.Sprintf("Deleted product ID: %s", productID)
	w.Write([]byte(response))
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")

	response := fmt.Sprintf("Search results - Query: %s, Min Price: %s, Max Price: %s", query, minPrice, maxPrice)
	w.Write([]byte(response))
}

func getProductsByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	sort := r.URL.Query().Get("sort")

	response := fmt.Sprintf("Products in category: %s, Sorted by: %s", category, sort)
	w.Write([]byte(response))
}

func getSimilarProducts(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	limit := r.URL.Query().Get("limit")

	response := fmt.Sprintf("Similar products to ID: %s, Limit: %s", productID, limit)
	w.Write([]byte(response))
}
