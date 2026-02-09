package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetTodaySales() ([]models.Transaction, error) {
	return s.repo.GetTodaySales()
}

func (s *ReportService) GetDailySales(start_date string, end_date string) ([]models.Transaction, error) {
	return s.repo.GetDailySales(start_date, end_date)
}
