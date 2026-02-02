package models

// Product adalah struct yang merepresentasikan data produk.
// Di PHP, ini mirip dengan Class atau Entity yang hanya berisi properti (Data Transfer Object).
type Product struct {
	// Field struct didefinisikan dengan Nama TipeData.
	// `json:"id"` adalah struct tag. Ini memberi tahu Go bagaimana field ini harus ditulis/dibaca dari JSON.
	// Jadi saat jadi JSON, key-nya akan menjadi "id", bukan "ID".
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}
