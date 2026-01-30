package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Goreng", Harga: 3500, Stok: 100},
	{ID: 2, Nama: "Le Minenar 600ml", Harga: 3000, Stok: 50},
	{ID: 3, Nama: "Ultra Milk 1L", Harga: 20000, Stok: 200},
}

// main adalah fungsi entry point aplikasi kasir API.
// Fungsi ini mengatur routing HTTP dan menjalankan server pada port 8080.
//
// Routes yang tersedia:
// - GET /health: Endpoint untuk pengecekan kesehatan API
// - GET /api/produk: Mengambil daftar semua produk
// - POST /api/produk: Membuat produk baru
func main() {

	//ROUTE Categoies

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

	//ROUTE Produk

	http.HandleFunc(
		"/api/produk/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				getProdukByID(w, r)
			} else if r.Method == "PUT" {
				updateProduk(w, r)
			} else if r.Method == "DELETE" {
				deleteProduk(w, r)
			}
		})

	// Handler untuk endpoint /api/produk
	// Mendukung operasi GET (read) dan POST (create)
	http.HandleFunc(
		"/api/produk", func(w http.ResponseWriter, r *http.Request) {
			// GET request: mengembalikan seluruh list produk
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(produk)
			} else if r.Method == "POST" {
				// POST request: membuat produk baru
				// Decode JSON dari request body ke struct Produk
				var produkBaru Produk
				err := json.NewDecoder(r.Body).Decode(&produkBaru)
				if err != nil {
					// Jika decode gagal, kirim error response
					http.Error(w, "Invalid request", http.StatusBadRequest)
				}

				// Generate ID otomatis berdasarkan jumlah produk yang ada
				produkBaru.ID = len(produk) + 1
				// Tambahkan produk baru ke slice produk
				produk = append(produk, produkBaru)

				// Set response header dan status code
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated) // HTTP 201 Created
				// Kirim produk yang baru dibuat sebagai response
				json.NewEncoder(w).Encode(produkBaru)
			}
		})

	// Handler untuk endpoint /health
	// Mengembalikan status API dalam format JSON
	// Response: {"status": "OK", "message": "API Running"}
	http.HandleFunc(
		"/health", func(w http.ResponseWriter, r *http.Request) {
			// Set header response sebagai JSON
			w.Header().Set("Content-Type", "application/json")
			// Encode dan kirim response JSON
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "OK",
				"message": "API Running",
			})
		})

	fmt.Println("Server running di localhost:8080")

	// Jalankan HTTP server pada port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		// Jika server gagal berjalan, tampilkan pesan error
		fmt.Println("Gagal Running Server")
	}
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/produk/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	//cari produk dengan id tersebut
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	//kalau not found
	http.Error(w, "Produk Belum Ada", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	//get produk id

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	//ganti int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	//get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// loop produk, cari id, ganti sesuai data dari request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)

			return
		}
	}

	http.Error(w, "Produk Belum Ada", http.StatusNotFound)

}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	//get produk ID

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	//convert id dari string ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			//buat slice baru dengan data sebelum dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "hapus berhasil",
			})
			return
		}
	}

	http.Error(w, "Produk Belum Ada", http.StatusNotFound)

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
