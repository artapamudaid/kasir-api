# Kasir API (Go Learning Project)

Project ini adalah RESTful API sederhana untuk manajemen produk Kasir, dibuat sebagai bahan belajar transisi dari **PHP** ke **Golang**. Project ini menerapkan konsep *Clean Architecture* sederhana.

## 🚀 Fitur

- **CRUD Produk**: Create, Read, Update, Delete data produk.
- **Database**: Menggunakan PostgreSQL (driver `lib/pq`).
- **Config Management**: Menggunakan `Viper` untuk handle environment variables (.env).
- **Architecture**: Layered pattern (Handler -> Service -> Repository -> Database).

## 🛠️ Tech Stack

- **Languange**: Go (Golang)
- **Database**: PostgreSQL (via Supabase/Local)
- **Router**: `net/http` (Standard Library) - Tanpa framework berat seperti Laravel!
- **Config**: `github.com/spf13/viper`

## 📂 Struktur Project (Analogi PHP)

Untuk mempermudah pemahaman programmer PHP:

| BuildGo Component | PHP Analogy (Laravel/CI) | Deskripsi |
|-------------------|--------------------------|-----------|
| `main.go` | `public/index.php` | Entry point aplikasi. Setup config & wiring. |
| `handlers/` | `Controllers/` | Menerima Request, validasi input, kirim Response (JSON). |
| `services/` | `Services/` / `Logic` | Business logic. Validasi data sebelum ke DB. |
| `repositories/` | `Models/` | Query SQL langsung. Tidak ada Eloquent di sini! |
| `models/` | `DTO` / Entity | Struct data. Class yang hanya berisi properti. |

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
   Isi `.env`:
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
   *Note: Saat pertama kali dijalankan, aplikasi akan otomatis membuat tabel `products` jika belum ada.*

## 🔌 API Endpoints

### Produk (Clean Architecture)

| Method | Endpoint | Deskripsi | Body Request (JSON) |
|--------|----------|-----------|---------------------|
| `GET` | `/api/produk` | List semua produk | - |
| `POST` | `/api/produk` | Tambah produk baru | `{"name": "...", "price": 1000, "stock": 10}` |
| `GET` | `/api/produk/{id}`| Detail produk | - |
| `POST` | `/api/produk/{id}`| Update produk | `{"name": "...", "price": 1000, "stock": 10}` |
| `DELETE`| `/api/produk/{id}`| Hapus produk | - |

*(Catatan: Untuk Update biasanya menggunakan PUT/PATCH, tapi di sini menggunakan POST untuk penyederhanaan belajar handling HTTP method di Go standard library).*

---
Happy Coding with Go! 🐹
