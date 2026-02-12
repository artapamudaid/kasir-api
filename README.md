# Kasir API (Go Learning Project)

Project ini adalah RESTful API sederhana untuk manajemen produk dan kategori Kasir, dibuat sebagai bahan belajar transisi dari **PHP** ke **Golang**. Project ini menerapkan konsep *Clean Architecture* (Handler -> Service -> Repository) dan **Relational Database**.

## 🚀 Fitur

- **CRUD Produk**: Create, Read, Update, Delete data produk (termasuk relasi ke Kategori).
- **CRUD Kategori**: Create, Read, Update, Delete data kategori.
- **Relasi Data**: Join Table antara Products dan Categories (One-to-Many).
- **Checkout Transaksi**: Endpoint untuk membuat transaksi belanja (mengurangi stok produk secara otomatis).
- **Auto Migration**: Tabel database (`products`, `categories`, `transactions`, `transaction_details`) dibuat otomatis saat aplikasi berjalan.
- **Database**: Menggunakan PostgreSQL (driver `lib/pq`).
- **Config Management**: Menggunakan `Viper` untuk handle environment variables (.env).
- **Security**: Middleware sederhana menggunakan Static API Key (`X-API-Key`).

## 🛠️ Tech Stack

- **Language**: Go (Golang) 1.25+
- **Database**: PostgreSQL (via Supabase/Local)
- **Router**: `net/http` (Standard Library) - Tanpa framework berat seperti Laravel!
- **Config**: `github.com/spf13/viper`

## 📂 Struktur Project (Analogi PHP)

Untuk mempermudah pemahaman programmer PHP:

| BuildGo Component | PHP Analogy (Laravel/CI) | Deskripsi |
|-------------------|--------------------------|-----------|
| `main.go` | `public/index.php` | Entry point aplikasi. Setup config, DB connection & wiring. |
| `handlers/` | `Controllers/` | Menerima Request, validasi input, kirim Response (JSON). |
| `services/` | `Services/` / `Logic` | Business logic. Validasi data sebelum ke DB. |
| `repositories/` | `Models/` | Query SQL langsung `(SELECT JOIN, INSERT, dll)`. |
| `models/` | `DTO` / Entity | Struct data. Class yang hanya berisi properti. |
| `database/` | `config/database.php` | Init connection & Migration script. |

## ⚙️ Cara Menjalankan

1. **Clone Repository**
   ```bash
   git clone <repo-url>
   cd kasir-api
   ```

2. **Setup Environment**
   Duplicate file `.env-example` ke `.env` dan sesuaikan isinya.
   ```bash
   cp .env-example .env
   ```
   Isi `.env` (Pastikan URL Database benar):
   ```ini
   PORT=8080
   DB_CONN=postgresql://user:password@host:port/dbname?sslmode=disable
   API_KEY=rahasia123
   ```

3. **Install Dependencies**
   Mirip `composer install`.
   ```bash
   go mod tidy
   ```

4. **Jalankan Aplikasi**
   Mirip `php artisan serve`.
   ```bash
   go run main.go
   ```
   *Note: Saat pertama kali dijalankan, aplikasi akan otomatis menjalankan MIGRATION (membuat tabel products & categories).*

## 🔌 API Endpoints
   
Semua request ke endpoint di bawah ini (kecuali `/health`) **WAJIB** menyertakan header:
`X-API-Key: <nilai_api_key_di_env>`

### 📦 Produk
Response produk sudah termasuk `category_name` (JOIN).

| Method | Endpoint | Deskripsi | Body Request (JSON) |
|--------|----------|-----------|---------------------|
| `GET` | `/api/product` | List semua produk | - |
| `POST` | `/api/product` | Tambah produk baru | `{"name": "...", "price": 1000, "stock": 10, "category_id": 1}` |
| `GET` | `/api/product/{id}`| Detail produk | - |
| `POST` | `/api/product/{id}`| Update produk | `{"name": "...", "price": 1000, "stock": 10, "category_id": 1}` |
| `DELETE`| `/api/product/{id}`| Hapus produk | - |

### 🏷️ Kategori

| Method | Endpoint | Deskripsi | Body Request (JSON) |
|--------|----------|-----------|---------------------|
| `GET` | `/api/category` | List semua kategori | - |
| `POST` | `/api/category` | Tambah kategori baru | `{"name": "...", "description": "..."}` |
| `GET` | `/api/category/{id}`| Detail kategori | - |
| `POST` | `/api/category/{id}`| Update kategori | `{"name": "...", "description": "..."}` |
| `DELETE`| `/api/category/{id}`| Hapus kategori | - |

### 🛒 Transaksi (Checkout)

| Method | Endpoint | Deskripsi | Body Request (JSON) |
|--------|----------|-----------|---------------------|
| `POST` | `/api/checkout` | Checkout belanjaan | `{"items": [{"product_id": 1, "quantity": 2}]}` |

### 📊 Laporan (Report)

| Method | Endpoint | Deskripsi | Query Parameters |
|--------|----------|-----------|------------------|
| `GET` | `/api/report/hari-ini` | Laporan penjualan hari ini | - |
| `GET` | `/api/report` | Laporan penjualan berdasarkan tanggal | `?start_date=2024-01-01&end_date=2024-12-31` |

*(Catatan: Update dan Delete menggunakan method POST & DELETE sesuai standar REST, meskipun untuk update di sini disederhanakan menggunakan POST endpoint yang sama dengan Detail jika diperlukan logic khusus, namun implementasi code menggunakan POST pada handler update).*

---
Happy Coding with Go! 🐹
