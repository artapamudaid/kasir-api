package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

// CategoryRepository bertanggung jawab untuk interaksi langsung dengan database (SQL).
// Ini mirip dengan layer "Model" atau "Repository" di framework PHP (seperti Laravel/CI).
type CategoryRepository struct {
	db *sql.DB // Menyimpan koneksi database agar bisa dipakai di semua method
}

// NewCategoryRepository adalah function "constructor".
// Di Go tidak ada class/constructor bawaan, jadi kita buat function yang mengembalikan instance struct.
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll mengambil semua data kategori dari database.
// (repo *CategoryRepository) artinya method ini milik struct CategoryRepository (seperti $this di PHP).
func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories"

	// Eksekusi query query. Mengembalikan *sql.Rows (cursor hasil query)
	rows, err := repo.db.Query(query)

	// Cek error standar di Go. Kalau err != nil, berarti ada masalah.
	if err != nil {
		return nil, err
	}

	// 'defer' memastikan rows.Close() dipanggil saat function ini selesai dieksekusi.
	// Sangat penting untuk menutup koneksi agar tidak memory leak.
	defer rows.Close()

	// Inisialisasi slice (array dinamis) kosong untuk menampung hasil
	categories := make([]models.Category, 0)

	// Loop selama masih ada baris data (rows.Next() mirip while($row = fetch...) di PHP)
	for rows.Next() {
		var p models.Category
		// Scan menyalin data kolom ke variabel goal. Urutannya harus sama dengan SELECT.
		// &p.ID artinya kita kirim alamat memori (pointer) variabel agar bisa diisi nilainya.
		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}

		// Tambahkan data ke slice
		categories = append(categories, p)
	}

	return categories, nil
}

// Create menyimpan kategori baru.
func (repo *CategoryRepository) Create(category *models.Category) error {
	// $1, $2, $3 adalah placeholder parameter (untuk mencegah SQL Injection).
	// RETURNING id digunakan (di PostgreSQL) untuk langsung dapat ID yang baru digenerate.
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"

	// QueryRow digunakan karena kita cuma mengharapkan 1 baris hasil (yaitu ID).
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}

// GetById mencari satu kategori berdasarkan ID.
func (repo *CategoryRepository) GetById(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var p models.Category
	// QueryRow untuk ambil 1 data.
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description)

	// Handle khusus jika data tidak ditemukan
	if err == sql.ErrNoRows {
		return nil, errors.New("kategori tidak ditemukan")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Update memperbarui data kategori.
func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"

	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	// Cek apakah ada baris yang benar-benar berubah
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return nil
}

// Delete menghapus kategori berdasarkan ID.
func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
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
		return errors.New("kategori tidak ditemukan")
	}

	return err

}
