package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

// NewTransactionService adalah constructor untuk Service.
// Menerima dependency Repository (Dependency Injection Pattern).
func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

// Checkout memproses pembelian item.
// Method ini bertindak sebagai orchestrator untuk transaksi.
func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}
