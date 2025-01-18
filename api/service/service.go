package service

import (
	"context"

	"github.com/marcelluseasley/receipt-processor/api/models"
	pp "github.com/marcelluseasley/receipt-processor/pointsprocessor"
	"github.com/marcelluseasley/receipt-processor/repo"

	"github.com/google/uuid"
)

type Service interface {
	ProcessReceipt(ctx context.Context, receipt models.Receipt) (models.ReceiptResponse, error)
	GetPoints(ctx context.Context, receiptId string) (*models.PointsResponse, error)
}

type ReceiptService struct {
	db repo.Repository
}

func NewReceiptService(db repo.Repository) Service {
	return ReceiptService{
		db,
	}
}

func (s ReceiptService) ProcessReceipt(ctx context.Context, receipt models.Receipt) (models.ReceiptResponse, error) {

	receiptId := uuid.New().String()
	s.db.SavePoints(ctx, receiptId, pp.ProcessPoints(receipt))
	return models.ReceiptResponse{Id: receiptId}, nil
}

func (s ReceiptService) GetPoints(ctx context.Context, receiptId string) (*models.PointsResponse, error) {

	points, err := s.db.GetPoints(ctx, receiptId)
	if err != nil {
		return nil, err
	}
	return &models.PointsResponse{Points: points}, nil
}
