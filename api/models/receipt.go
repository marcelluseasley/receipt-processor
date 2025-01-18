package models

import (
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type ReceiptRequest struct {
	Retailer     string `json:"retailer" validate:"required,receipt_retailer"`
	PurchaseDate string `json:"purchaseDate" validate:"required,receipt_date"` // 2022-01-01
	PurchaseTime string `json:"purchaseTime" validate:"required,receipt_time"` // 13:01
	Total        string `json:"total" validate:"required,receipt_price"`
	Items        []struct {
		ShortDescription string `json:"shortDescription" validate:"required,receipt_shortdesc"`
		Price            string `json:"price" validate:"required,receipt_price"`
	} `json:"items" validate:"required"`
}

func (rr *ReceiptRequest) ToReceipt() (*Receipt, error) {
	r := &Receipt{}

	purchaseDate, err := time.Parse(time.DateOnly, rr.PurchaseDate)
	if err != nil {
		return nil, err
	}

	purchaseTime, err := time.Parse("15:04", rr.PurchaseTime)
	if err != nil {
		return nil, err
	}

	total, err := strconv.ParseFloat(rr.Total, 64)
	if err != nil {
		return nil, err
	}

	items := []Item{}

	for _, item := range rr.Items {
		i := Item{}

		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return nil, err
		}
		i.Price = price
		i.ShortDescription = item.ShortDescription
		items = append(items, i)
	}

	r.Retailer = rr.Retailer
	r.PurchaseDate = purchaseDate
	r.PurchaseTime = purchaseTime
	r.Total = total
	r.Items = items

	return r, nil
}

type Receipt struct {
	Retailer     string
	PurchaseDate time.Time
	PurchaseTime time.Time
	Total        float64
	Items        []Item
}

type Item struct {
	ShortDescription string
	Price            float64
}

type ReceiptResponse struct {
	Id string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

// validations based on api.yml

var retailerRegex = regexp.MustCompile(`^[\w\s\-&]+$`)
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var timeRegex = regexp.MustCompile(`^\d{2}:\d{2}$`)
var shortDescRegex = regexp.MustCompile(`^[\w\s\-&]+$`)
var priceRegex = regexp.MustCompile(`^\d+\.\d{2}$`)

func ValidateRetailer(fl validator.FieldLevel) bool {
	return retailerRegex.MatchString(fl.Field().String())
}

func ValidateDate(fl validator.FieldLevel) bool {
	return dateRegex.MatchString(fl.Field().String())
}

func ValidateTime(fl validator.FieldLevel) bool {
	return timeRegex.MatchString(fl.Field().String())
}

func ValidatePrice(fl validator.FieldLevel) bool {
	return priceRegex.MatchString(fl.Field().String())
}

func ValidateShortDescription(fl validator.FieldLevel) bool {
	return shortDescRegex.MatchString(fl.Field().String())
}
