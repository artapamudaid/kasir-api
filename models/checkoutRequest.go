package models

// CheckoutRequest adalah struct untuk payload JSON saat checkout.
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}
