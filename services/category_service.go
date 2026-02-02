package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

// CategoryService bertanggung jawab atas business logic.
// Di sinilah validasi data, kalkulasi harga, atau aturan bisnis lainnya dilakukan sebelum masuk ke database.
// Layer ini dipanggil oleh Handler (Controller) dan memanggil Repository.
type CategoryService struct {
	repo *repositories.CategoryRepository
}

// NewCategoryService adalah constructor untuk Service.
// Menerima dependency Repository (Dependency Injection Pattern).
func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll mengambil semua kategori, saat ini hanya meneruskan ke repository.
func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

// Create membuat kategori baru.
// Di real application, kita bisa tambahkan validasi di sini (misal: harga tidak boleh negatif).
func (s *CategoryService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

func (s *CategoryService) GetById(id int) (*models.Category, error) {
	return s.repo.GetById(id)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
