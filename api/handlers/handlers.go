package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/marcelluseasley/receipt-processor/api/models"
	"github.com/marcelluseasley/receipt-processor/api/service"
)

type Handler struct {
	service  service.Service
	validate *validator.Validate
}

func NewHandler(service service.Service, validate *validator.Validate) *Handler {
	return &Handler{
		service,
		validate,
	}
}

func (h Handler) PointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	receiptId := chi.URLParam(r, "id")

	pointsResp, err := h.service.GetPoints(ctx, receiptId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(pointsResp)
	if err != nil {
		handleInvalidReceipt(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h Handler) ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var receiptRequest models.ReceiptRequest

	err := json.NewDecoder(r.Body).Decode(&receiptRequest)
	if err != nil {
		handleInvalidReceipt(w)
		return
	}

	err = h.validate.Struct(receiptRequest)
	if err != nil {
		handleInvalidReceipt(w)
		return
	}

	receipt, err := receiptRequest.ToReceipt()
	if err != nil {
		handleInvalidReceipt(w)
		return
	}

	receiptResp, err := h.service.ProcessReceipt(ctx, *receipt)
	if err != nil {
		handleInvalidReceipt(w)
		return
	}

	jsonResponse, err := json.Marshal(receiptResp)
	if err != nil {
		handleInvalidReceipt(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleInvalidReceipt(w http.ResponseWriter) {
	http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
}
