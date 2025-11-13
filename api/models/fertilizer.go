package models

import "time"

type Fertilizer struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ProductName       string    `json:"productName"`
	Category          string    `json:"category"`
	Dosage            string    `json:"dosage"`
	UniqueID          string    `json:"uniqueID"`
	BatchNumber       string    `json:"batchNumber"`
	ManufactureDate   time.Time `json:"manufactureDate"`
	ExpiryDate        time.Time `json:"expiryDate"`
	CautionaryLogo    string    `json:"cautionaryLogo"`
	CautionaryText    string    `json:"cautionaryText"`
	AntidoteStatement string    `json:"antidoteStatement"`
	CreatedAt         time.Time `json:"createdAt"`
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
