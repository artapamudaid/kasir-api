# Kasir API (Go Learning Project)

Project ini adalah RESTful API sederhana untuk manajemen produk dan kategori Kasir, dibuat sebagai bahan belajar transisi dari **PHP** ke **Golang**. Project ini menerapkan konsep *Clean Architecture* (Handler -> Service -> Repository).

## 🚀 Fitur

- **CRUD Produk**: Create, Read, Update, Delete data produk.
- **CRUD Kategori**: Create, Read, Update, Delete data kategori.
- **Auto Migration**: Tabel database (`products`, `categories`) dibuat otomatis saat aplikasi berjalan.
- **Database**: Menggunakan PostgreSQL (driver `lib/pq`).
- **Config Management**: Menggunakan `Viper` untuk handle environment variables (.env).

## 🛠️ Tech Stack

- **Language**: Go (Golang) 1.21+
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
| `repositories/` | `Models/` | Query SQL langsung `(SELECT, INSERT, dll)`. |
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

### 📦 Produk

| Method | Endpoint | Deskripsi | Body Request (JSON) |
|--------|----------|-----------|---------------------|
| `GET` | `/api/produk` | List semua produk | - |
| `POST` | `/api/produk` | Tambah produk baru | `{"name": "...", "price": 1000, "stock": 10}` |
| `GET` | `/api/produk/{id}`| Detail produk | - |
| `POST` | `/api/produk/{id}`| Update produk | `{"name": "...", "price": 1000, "stock": 10}` |
| `DELETE`| `/api/produk/{id}`| Hapus produk | - |

### 🏷️ Kategori

| Method | Endpoint | Deskripsi | Body Request (JSON) |
|--------|----------|-----------|---------------------|
| `GET` | `/api/categories` | List semua kategori | - |
| `POST` | `/api/categories` | Tambah kategori baru | `{"name": "...", "description": "..."}` |
| `GET` | `/api/categories/{id}`| Detail kategori | - |
| `POST` | `/api/categories/{id}`| Update kategori | `{"name": "...", "description": "..."}` |
| `DELETE`| `/api/categories/{id}`| Hapus kategori | - |

*(Catatan: Update dan Delete menggunakan method POST & DELETE sesuai standar REST, meskipun untuk update di sini disederhanakan menggunakan POST endpoint yang sama dengan Detail jika diperlukan logic khusus, namun implementasi code menggunakan POST pada handler update).*

---
Happy Coding with Go! 🐹
