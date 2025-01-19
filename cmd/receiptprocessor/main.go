package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

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

	logger := logrus.New()

	/*
		Since the api.yml file only lists 2 types of errors,
		we need to at least have logs show where the actual
		problem occurs...hence, SetReportCaller
	*/
	logger.SetReportCaller((true))

	fmt.Printf("server started on port %s", port)
	logger.Fatal(newServer(port, logger).ListenAndServe())

}

func newServer(port string, logger *logrus.Logger) *Server {
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
	h := handlers.NewHandler(rs, validate, logger)

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
