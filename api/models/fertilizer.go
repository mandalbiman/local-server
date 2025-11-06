package models

import "time"

type Fertilizer struct {
	ID                uint `gorm:"primaryKey"`
	ProductName       string
	Category          string
	Dosage            string
	UniqueID          string
	BatchNumber       string
	ManufactureDate   time.Time
	ExpiryDate        time.Time
	CautionaryLogo    string
	CautionaryText    string
	AntidoteStatement string
	CreatedAt         time.Time
}

// Used only for JSON import
type FertilizerJSON struct {
	ProductName       string `json:"product_name"`
	Category          string `json:"category"`
	Dosage            string `json:"dosage"`
	UniqueID          string `json:"unique_id"`
	BatchNumber       string `json:"batch_number"`
	ManufactureDate   string `json:"manufacture_date"`
	ExpiryDate        string `json:"expiry_date"`
	CautionaryLogo    string `json:"cautionary_logo"`
	CautionaryText    string `json:"cautionary_text"`
	AntidoteStatement string `json:"antidote_statement"`
}
