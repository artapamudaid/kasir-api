package middleware

import (
	"net/http"
)

// ApiKey adalah middleware untuk validasi API Key
// Mengembalikan fungsi yang menerima http.HandlerFunc dan mengembalikan http.HandlerFunc
func ApiKey(validApiKey string) func(http.HandlerFunc) http.HandlerFunc {
	// Fungsi ini akan menerima handler asli (next)
	return func(next http.HandlerFunc) http.HandlerFunc {
		// Fungsi ini akan menerima request dan response
		return func(w http.ResponseWriter, r *http.Request) {
			// Ambil API Key dari header
			apiKey := r.Header.Get("X-API-Key")

			// Cek apakah API Key ada
			if apiKey == "" {
				http.Error(w, "API Key is Required", http.StatusUnauthorized)
				return
			}

			// Cek apakah API Key valid
			if apiKey != validApiKey {
				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
				return
			}

			// API Key Valid, Lanjutkan ke Handler
			next(w, r)
		}
	}
}
