package models

import (
	"time"
)

type Status string
type PaymentType string

const (
	Accepted Status = "pending"
	Declined Status = "failed"
)

const (
	Cache PaymentType = "cache"
)

type Transaction struct {
	TransactionID      uint      `csv:"Transaction ID"`
	RequestID          uint      `csv:"Request ID"`
	TerminalID         uint      `csv:"Terminal ID"`
	PartnerObjectID    uint      `csv:"Partner Object ID"`
	AmountTotal        float64   `csv:"Amount Total"`
	AmountOriginal     float64   `csv:"Amount Original"`
	CommissionPS       float64   `csv:"Commission PS"`
	CommissionClient   float64   `csv:"Commission Client"`
	CommissionProvider float64   `csv:"Commission Provider"`
	DateInput          time.Time `csv:"Date Input"`
	DatePost           time.Time `csv:"Date Post"`
	Status             string    `csv:"Status"`
	PaymentType        string    `csv:"Payment Type"`
	PaymentNumber      string    `csv:"Payment Number"`
	ServiceID          uint      `csv:"Service ID"`
	Service            string    `csv:"Service"`
	PayeeID            uint      `csv:"Payee ID"`
	PayeeName          string    `csv:"Payee Name"`
	PayeeBankMfo       uint      `csv:"Payee Bank MFO"`
	PayeeBankAccount   string    `csv:"Payee Bank Account"`
	PaymentNarrative   string    `csv:"Payment Narrative"`    
}

type SuccessResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}