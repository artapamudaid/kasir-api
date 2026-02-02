package models

// Category adalah struct yang merepresentasikan data kategori.
// Di PHP, ini mirip dengan Class atau Entity yang hanya berisi properti (Data Transfer Object).
type Category struct {
	// Field struct didefinisikan dengan Nama TipeData.
	// `json:"id"` adalah struct tag. Ini memberi tahu Go bagaimana field ini harus ditulis/dibaca dari JSON.
	// Jadi saat jadi JSON, key-nya akan menjadi "id", bukan "ID".
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
