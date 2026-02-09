package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetTodaySales(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.service.GetTodaySales()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	start_date := r.URL.Query().Get("start_date")
	end_date := r.URL.Query().Get("end_date")

	transactions, err := h.service.GetDailySales(start_date, end_date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
