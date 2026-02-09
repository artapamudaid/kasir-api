package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

// ProductService bertanggung jawab atas business logic.
// Di sinilah validasi data, kalkulasi harga, atau aturan bisnis lainnya dilakukan sebelum masuk ke database.
// Layer ini dipanggil oleh Handler (Controller) dan memanggil Repository.
type ProductService struct {
	repo *repositories.ProductRepository
}

// NewProductService adalah constructor untuk Service.
// Menerima dependency Repository (Dependency Injection Pattern).
func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAll mengambil semua produk, saat ini hanya meneruskan ke repository.
func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

// Create membuat produk baru.
// Di real application, kita bisa tambahkan validasi di sini (misal: harga tidak boleh negatif).
func (s *ProductService) Create(data *models.Product) error {
	return s.repo.Create(data)
}

// GetById mengambil satu produk secara spesifik.
func (s *ProductService) GetById(id int) (*models.Product, error) {
	return s.repo.GetById(id)
}

// Update memvalidasi dan memperbarui data produk.
func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

// Delete menghapus produk (soft delete atau hard delete, tergantung implementasi repo).
func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
