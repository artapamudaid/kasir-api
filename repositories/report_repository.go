package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetTodaySales() ([]models.Transaction, error) {
	query := "SELECT id, total_amount, created_at FROM transactions WHERE DATE(created_at) = $1"
	date := time.Now().Format("2006-01-02")

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *ReportRepository) GetDailySales(start_date string, end_date string) ([]models.Transaction, error) {
	query := "SELECT id, total_amount, created_at FROM transactions"

	args := []interface{}{}

	if start_date != "" && end_date != "" {
		query += " WHERE DATE(created_at) BETWEEN $1 AND $2"
		args = append(args, start_date, end_date)
	} else {
		return nil, errors.New("start_date and end_date are required")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
