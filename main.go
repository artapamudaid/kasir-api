package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middleware"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config struct untuk menampung konfigurasi dari environment variables atau file .env
type Config struct {
	Port    string `mapstructure:"PORT"`    // Ambil dari key PORT
	DBConn  string `mapstructure:"DB_CONN"` // Ambil dari key DB_CONN
	API_KEY string `mapstructure:"API_KEY"` // Ambil dari key API_KEY
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
		Port:    viper.GetString("PORT"),
		DBConn:  viper.GetString("DB_CONN"),
		API_KEY: viper.GetString("API_KEY"),
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
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	reportRepo := repositories.NewReportRepository(db)
	// 2. Service butuh Repository
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	reportService := services.NewReportService(reportRepo)
	// 3. Handler (Controller) butuh Service
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)

	// Setup Middleware
	apiKeyMiddleware := middleware.ApiKey(config.API_KEY)

	// === ROUTING ===
	// Menggunakan Default ServeMux dari package net/http

	// Route untuk PRODUK (Menggunakan Pattern Layered Architecture yang proper)
	http.HandleFunc("/api/produk", apiKeyMiddleware(productHandler.HandleProducts))
	http.HandleFunc("/api/produk/", apiKeyMiddleware(productHandler.HandleProductByID))

	// Route untuk CATEGORIES (Menggunakan Pattern Layered Architecture yang proper)
	http.HandleFunc("/api/categories", apiKeyMiddleware(categoryHandler.HandleCategories))
	http.HandleFunc("/api/categories/", apiKeyMiddleware(categoryHandler.HandleCategoryByID))

	// Route untuk TRANSACTIONS (Menggunakan Pattern Layered Architecture yang proper)
	http.HandleFunc("/api/checkout", apiKeyMiddleware(transactionHandler.Checkout))

	// Route untuk REPORT
	http.HandleFunc("/api/report/hari-ini", apiKeyMiddleware(reportHandler.GetTodaySales))
	http.HandleFunc("/api/report", apiKeyMiddleware(reportHandler.GetReport))

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
