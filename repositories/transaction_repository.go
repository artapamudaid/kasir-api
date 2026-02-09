package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

// TransactionRepository menangani transaksi database terkait order/penjualan.
type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction membuat transaksi baru dengan atomic transaction (ACID).
// Ini memastikan stock berkurang dan transaksi tercatat, atau tidak sama sekali jika terjadi error.
func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := repo.db.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Bulk Insert Optimization
	// Menggabungkan semua insert menjadi satu query untuk mengurangi Round-Trip Time (RTT) ke database.
	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		vals := []interface{}{}

		for i, detail := range details {
			// Update struct di memori agar response lengkap
			details[i].TransactionID = transactionID

			// Buat placeholder ($1, $2, $3, $4), ($5, $6, $7, $8), dst...
			n := i * 4
			query += fmt.Sprintf("($%d, $%d, $%d, $%d),", n+1, n+2, n+3, n+4)
			vals = append(vals, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		}

		// Hapus koma terakhir
		query = query[:len(query)-1]

		// Eksekusi satu query besar
		_, err = tx.Exec(query, vals...)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
