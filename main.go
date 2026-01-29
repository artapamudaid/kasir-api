package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func main() {
	// localhost:8080/health
	http.HandleFunc(
		"/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "OK",
				"message": "API Running",
			})
		})

	http.HandleFunc(
		"/api/produk", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(produk)
			} else if r.Method == "POST" {
				var produkBaru Produk
				err := json.NewDecoder(r.Body).Decode(&produkBaru)
				if err != nil {
					http.Error(w, "Invalid request", http.StatusBadRequest)
				}

				produkBaru.ID = len(produk) + 1
				produk = append(produk, produkBaru)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(produkBaru)
			}
		})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gagal Running Server")
	} else {
		fmt.Println("Server Running on PORT 8080")
	}

}
