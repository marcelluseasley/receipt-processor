package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

	"github.com/marcelluseasley/receipt-processor/api/handlers"
	"github.com/marcelluseasley/receipt-processor/api/models"
	"github.com/marcelluseasley/receipt-processor/api/service"
	"github.com/marcelluseasley/receipt-processor/repo"
)

type Server struct {
	*http.Server
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Infof("server started on port %s", port)
	log.Fatal(newServer(port).ListenAndServe())

}

func newServer(port string) *Server {
	router := chi.NewRouter()

	// custom validations for the regex requirements (api.yml)
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("receipt_retailer", models.ValidateRetailer)
	validate.RegisterValidation("receipt_shortdesc", models.ValidateShortDescription)
	validate.RegisterValidation("receipt_price", models.ValidatePrice)
	validate.RegisterValidation("receipt_date", models.ValidateDate)
	validate.RegisterValidation("receipt_time", models.ValidateTime)

	// set up dependencies
	db := repo.NewRepository()
	rs := service.NewReceiptService(db)
	h := handlers.NewHandler(rs, validate)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	server := &Server{
		httpServer,
	}

	router.Get("/receipts/{id}/points", h.PointsHandler)
	router.Post("/receipts/process", h.ProcessReceiptHandler)

	return server
}
