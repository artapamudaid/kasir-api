package models

// CheckoutItem merepresentasikan satu item produk dalam keranjang belanja.
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
