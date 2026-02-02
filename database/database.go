package database

import (
	"database/sql"
	"log"

	// Import driver PostgreSQL.
	// Tanda underscores (_) berarti kita hanya menjalankan fungsi init() dari library ini (side-effect),
	// tanpa menggunakan function-nya secara langsung.
	// Ini mendaftarkan driver "postgres" ke package database/sql.
	// Mirip seperti mengaktifkan extension pgsql di php.ini.
	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	//open database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	//test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	//set connection pool settings (optional tapi recomended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connection successfully")
	return db, nil
}
