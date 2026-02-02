package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

// ProductHandler mirip dengan Controller di MVC PHP (Laravel/CI).
// Tugasnya: Menerima Request HTTP -> Validasi Input Dasar -> Panggil Service -> Kembalikan Response JSON.
type ProductHandler struct {
	service *services.ProductService
}

// NewProductHandler adalah constructor.
func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

//SELECT & CREATE HANDLER

// HandleProducts adalah endpoint untuk /api/produk (list & create).
// Menerapkan REST API standard dimana GET untuk baca dan POST untuk tulis.
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	// r.Method berisi tipe method HTTP (GET, POST, dll).
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		// Jika method tidak didukung, kembalikan 405 Method Not Allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProductByID adalah endpoint untuk /api/produk/{id} (detail, update, delete).
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetById(w, r)
	case http.MethodPost: // Note: Biasanya Update menggunakan PUT atau PATCH di RESTful API yang ketat.
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAll mengambil semua data dan mengirimkan JSON response.
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Panggil service untuk ambil data
	products, err := h.service.GetAll()
	if err != nil {
		// http.Error adalah helper untuk mengirim error response sederhana
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type agar client tahu ini JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode data struct ke JSON dan tulis ke ResponseWriter
	json.NewEncoder(w).Encode(products)
}

// Create membaca JSON dari request body dan menyimpan data.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	// Decode JSON dari body request ke struct product
	// Ini mirip json_decode(file_get_contents('php://input')) di PHP Native
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated) // HTTP 201 Created
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	// Parsing ID dari URL manual, karena pakai standard library net/http (agak ribet dibanding framework seperti Gin/Echo).
	// Mengambil string setelah "/api/produk/"
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr) // Convert string ke int

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	})
}
