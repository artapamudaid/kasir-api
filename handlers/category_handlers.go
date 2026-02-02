package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

// CategoryHandler mirip dengan Controller di MVC PHP (Laravel/CI).
// Tugasnya: Menerima Request HTTP -> Validasi Input Dasar -> Panggil Service -> Kembalikan Response JSON.
type CategoryHandler struct {
	service *services.CategoryService
}

// NewCategoryHandler adalah constructor.
// Menerima dependency Service yang dibutuhkan (Dependency Injection).
func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

//SELECT & CREATE HANDLER

// HandleCategories adalah endpoint untuk /api/categories (List & Create).
// Menerapkan REST API standard dimana GET untuk membaca data dan POST untuk menambah data.
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	// r.Method berisi tipe method HTTP yang dikirim client (GET, POST, dll).
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		// Jika method bukan GET atau POST, kembalikan 405 Method Not Allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleCategoryByID adalah endpoint untuk /api/categories/{id} (Detail, Update, Delete).
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetById(w, r)
	case http.MethodPost: // Note: Idealnya menggunakan PUT/PATCH untuk update, tapi POST juga bisa digunakan.
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAll mengambil semua data kategori dan mengirimkan JSON response.
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Panggil service business logic untuk ambil data
	categories, err := h.service.GetAll()
	if err != nil {
		// http.Error helper untuk kirim error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type agar client (browser/postman) tahu ini data JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode data dari struct Go menjadi JSON string dan tulis ke ResponseWriter
	json.NewEncoder(w).Encode(categories)
}

// Create membaca input JSON, validasi, dan menyimpan data kategori baru.
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	// Decode JSON dari body request ke struct category.
	// Ini mirip dengan json_decode(file_get_contents('php://input'), true) di PHP.
	err := json.NewDecoder(r.Body).Decode(&category)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil service untuk logic penyimpanan
	err = h.service.Create(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Kembalikan response sukses 201 Created dengan data yang baru dibuat
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	// Parsing ID dari URL manual, karena pakai standard library net/http (agak ribet dibanding framework seperti Gin/Echo).
	// Mengambil string setelah "/api/categories/"
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr) // Convert string ke int

	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category.ID = id
	err = h.service.Update(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Category deleted successfully",
	})
}
