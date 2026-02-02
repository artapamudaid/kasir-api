package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config struct untuk menampung konfigurasi dari environment variables atau file .env
type Config struct {
	Port   string `mapstructure:"PORT"`    // Ambil dari key PORT
	DBConn string `mapstructure:"DB_CONN"` // Ambil dari key DB_CONN
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var category = []Category{
	{ID: 1, Name: "Sembako", Description: "Ini sembako"},
	{ID: 2, Name: "Air Mineral", Description: "Ini air mineral"},
	{ID: 3, Name: "Susu", Description: "Ini susu"},
}

// main adalah Fungsi Utama (Entry Point) aplikasi Go. Mirip index.php di projecr PHP.
func main() {
	// Konfigurasi Viper untuk membaca Environment Variables & .env file
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Cek apakah ada file .env, jika ada maka load (untuk development lokal)
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// Load config ke struct
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database connection
	// Memanggil fungsi InitDB dari package database (buatan sendiri)
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	// Defer close untuk menutup koneksi jika main selesai (saat aplikasi mati)
	defer db.Close()

	// === WIRING DEPENDENCY INJECTION ===
	// Menghubungkan layer-layer aplikasi: Repository -> Service -> Handler
	// Ini dilakukan manual di Go (tanpa framework magie seperti Laravel Service Container)

	// 1. Repository butuh DB
	productRepo := repositories.NewProductRepository(db)
	// 2. Service butuh Repository
	productService := services.NewProductService(productRepo)
	// 3. Handler (Controller) butuh Service
	productHandler := handlers.NewProductHandler(productService)

	// === ROUTING ===
	// Menggunakan Default ServeMux dari package net/http

	// Route untuk PRODUK (Menggunakan Pattern Layered Architecture yang proper)
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID) // Typo di 'proud'? Harusnya 'product' mungkin :)

	// Note: Di bawah ini adalah contoh Route untuk CATEGORIES
	// Ini menggunakan gaya "Procedural/Flat" (tanpa layer Service/Repo terpisah).
	// Bagus untuk belajar perbandingan, tapi tidak disarankan untuk aplikasi besar.
	http.HandleFunc(
		"/api/categories/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				getCategoryByID(w, r)
			} else if r.Method == "PUT" {
				updateCategory(w, r)
			} else if r.Method == "DELETE" {
				deleteCategory(w, r)
			}
		})

	// Handler untuk endpoint /api/categories
	// Mendukung operasi GET (read) dan POST (create)
	http.HandleFunc(
		"/api/categories", func(w http.ResponseWriter, r *http.Request) {
			// GET request: mengembalikan seluruh list categories
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(category)
			} else if r.Method == "POST" {
				// POST request: membuat categories baru
				// Decode JSON dari request body ke struct Category
				var categoryBaru Category
				err := json.NewDecoder(r.Body).Decode(&categoryBaru)
				if err != nil {
					// Jika decode gagal, kirim error response
					http.Error(w, "Invalid request", http.StatusBadRequest)
				}

				// Generate ID otomatis berdasarkan jumlah category yang ada
				categoryBaru.ID = len(category) + 1
				// Tambahkan category baru ke slice category
				category = append(category, categoryBaru)

				// Set response header dan status code
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated) // HTTP 201 Created
				// Kirim category yang baru dibuat sebagai response
				json.NewEncoder(w).Encode(categoryBaru)
			}
		})

	// Handler untuk endpoint /health
	http.HandleFunc(
		"/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "OK",
				"message": "API Running",
			})
		})

	// Jalankan Server
	addr := "0.0.0.0:" + config.Port
	if config.Port == "" {
		addr = "0.0.0.0:8080"
	}
	fmt.Println("Server running in ", addr)

	// ListenAndServe memblokir proses main dan terus mendengarkan request masuk
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Failed to run server", err)
	}
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/categories/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	//cari category dengan id tersebut
	for _, p := range category {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	//kalau not found
	http.Error(w, "Category Belum Ada", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	//get category id

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	//ganti int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	//get data dari request
	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// loop category, cari id, ganti sesuai data dari request
	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)

			return
		}
	}

	http.Error(w, "Category Belum Ada", http.StatusNotFound)

}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	//get category ID

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	//convert id dari string ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	for i, p := range category {
		if p.ID == id {
			//buat slice baru dengan data sebelum dan sesudah index
			category = append(category[:i], category[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "hapus berhasil",
			})
			return
		}
	}

	http.Error(w, "Category Belum Ada", http.StatusNotFound)

}
