package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/marcelluseasley/receipt-processor/api/models"
	"github.com/marcelluseasley/receipt-processor/api/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service  service.Service
	validate *validator.Validate
	logger   *logrus.Logger
}

func NewHandler(service service.Service, validate *validator.Validate, logger *logrus.Logger) *Handler {
	return &Handler{
		service,
		validate,
		logger,
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
		h.logger.Errorf("unable to marshal response: %v", err)
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
		h.logger.Errorf("unable to decode body into receiptRequest struct: %v", err)
		handleInvalidReceipt(w)
		return
	}

	err = h.validate.Struct(receiptRequest)
	if err != nil {
		h.logger.Errorf("unable to validate receiptRequest struct: %v", err)
		handleInvalidReceipt(w)
		return
	}

	receipt, err := receiptRequest.ToReceipt()
	if err != nil {
		handleInvalidReceipt(w)
		return
	}

	receiptResp := h.service.ProcessReceipt(ctx, *receipt)

	jsonResponse, err := json.Marshal(receiptResp)
	if err != nil {
		h.logger.Errorf("unable to marshal receiptResp struct: %v", err)
		handleInvalidReceipt(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleInvalidReceipt(w http.ResponseWriter) {
	http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
}
