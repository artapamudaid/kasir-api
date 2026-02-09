package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

// ProductRepository bertanggung jawab untuk interaksi langsung dengan database (SQL).
// Ini mirip dengan layer "Model" atau "Repository" di framework PHP (seperti Laravel/CI).
type ProductRepository struct {
	db *sql.DB // Menyimpan koneksi database agar bisa dipakai di semua method
}

// NewProductRepository adalah function "constructor".
// Di Go tidak ada class/constructor bawaan, jadi kita buat function yang mengembalikan instance struct.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll mengambil semua data produk dari database.
// (repo *ProductRepository) artinya method ini milik struct ProductRepository (seperti $this di PHP).
func (repo *ProductRepository) GetAll(filterName string) ([]models.Product, error) {

	query := `
		SELECT 
			p.id, 
			p.name, 
			p.price, 
			p.stock, 
			c.id as category_id,
			c.name as category_name
		FROM products p
		JOIN categories c ON p.category_id = c.id`

	args := []interface{}{}
	if filterName != "" {
		query += "WHERE p.name LIKE $1"
		args = append(args, "%"+filterName+"%")
	}

	// Eksekusi query query. Mengembalikan *sql.Rows (cursor hasil query)
	rows, err := repo.db.Query(query, args...)

	// Cek error standar di Go. Kalau err != nil, berarti ada masalah.
	if err != nil {
		return nil, err
	}

	// 'defer' memastikan rows.Close() dipanggil saat function ini selesai dieksekusi.
	// Sangat penting untuk menutup koneksi agar tidak memory leak.
	defer rows.Close()

	// Inisialisasi slice (array dinamis) kosong untuk menampung hasil
	products := make([]models.Product, 0)

	// Loop selama masih ada baris data (rows.Next() mirip while($row = fetch...) di PHP)
	for rows.Next() {
		var p models.Product
		// Scan menyalin data kolom ke variabel goal. Urutannya harus sama dengan SELECT.
		// &p.ID artinya kita kirim alamat memori (pointer) variabel agar bisa diisi nilainya.
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
		if err != nil {
			return nil, err
		}

		// Tambahkan data ke slice
		products = append(products, p)
	}

	return products, nil
}

// Create menyimpan produk baru.
func (repo *ProductRepository) Create(product *models.Product) error {
	// $1, $2, $3 adalah placeholder parameter (untuk mencegah SQL Injection).
	// RETURNING id digunakan (di PostgreSQL) untuk langsung dapat ID yang baru digenerate.
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"

	// QueryRow digunakan karena kita cuma mengharapkan 1 baris hasil (yaitu ID).
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)

	if err != nil {
		return err
	}

	// Fix: Ambil nama kategori agar response JSON lengkap (tidak kosong string-nya)
	queryCategory := "SELECT name FROM categories WHERE id = $1"
	err = repo.db.QueryRow(queryCategory, product.CategoryID).Scan(&product.CategoryName)

	return err
}

// GetById mencari satu produk berdasarkan ID.
func (repo *ProductRepository) GetById(id int) (*models.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, 
				c.id as category_id, c.name as category_name
				FROM products p 
				JOIN categories c ON p.category_id = c.id 
				WHERE p.id = $1`

	var p models.Product
	// QueryRow untuk ambil 1 data.
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)

	// Handle khusus jika data tidak ditemukan
	if err == sql.ErrNoRows {
		return nil, errors.New("Produk tidak ditemukan")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Update memperbarui data produk.
func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"

	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	// Cek apakah ada baris yang benar-benar berubah
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

// Delete menghapus produk berdasarkan ID.
func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	// Exec digunakan untuk query yang tidak mengembalikan baris data (INSERT, UPDATE, DELETE)
	result, err := repo.db.Exec(query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err

}
