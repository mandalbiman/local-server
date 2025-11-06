package db

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"bats.com/local-server/api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FertilizerDB struct {
	db *gorm.DB
}

func GetFertilizerDB() *FertilizerDB {
	db := setUpFertilizerDB()
	db.AutoMigrate(&models.Fertilizer{}) // ensure table exists
	return &FertilizerDB{db: db}
}

func setUpFertilizerDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("fertilizers.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}

// Get all fertilizers
func (f FertilizerDB) GetFertilizers() ([]models.Fertilizer, error) {
	var list []models.Fertilizer
	result := f.db.Find(&list)
	return list, result.Error
}

// Get fertilizer by ID
func (f FertilizerDB) GetFertilizerByID(id int) (*models.Fertilizer, error) {
	var fert models.Fertilizer
	result := f.db.First(&fert, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fert, nil
}

// Create a fertilizer
func (f FertilizerDB) CreateFertilizer(fert models.Fertilizer) (*models.Fertilizer, error) {
	result := f.db.Create(&fert)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fert, nil
}

// Update fertilizer
func (f FertilizerDB) UpdateFertilizer(fert models.Fertilizer) (*models.Fertilizer, error) {
	result := f.db.Save(&fert)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fert, nil
}

// Delete fertilizer
func (f FertilizerDB) DeleteFertilizer(id int) (*models.Fertilizer, error) {
	fert, err := f.GetFertilizerByID(id)
	if err != nil {
		return nil, err
	}
	result := f.db.Delete(&fert)
	return fert, result.Error
}

// Search by product name
func (f FertilizerDB) SearchByProductName(keyword string) ([]models.Fertilizer, error) {
	var items []models.Fertilizer
	result := f.db.Where("product_name LIKE ?", "%"+keyword+"%").Find(&items)
	return items, result.Error
}

// Load fertilizers from JSON
func (f FertilizerDB) LoadFromJSON(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read JSON: %w", err)
	}

	var list []models.FertilizerJSON
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, item := range list {
		mDate, _ := time.Parse("02-01-2006", item.ManufactureDate)
		eDate, _ := time.Parse("02-01-2006", item.ExpiryDate)

		fert := models.Fertilizer{
			ProductName:       item.ProductName,
			Category:          item.Category,
			Dosage:            item.Dosage,
			UniqueID:          item.UniqueID,
			BatchNumber:       item.BatchNumber,
			ManufactureDate:   mDate,
			ExpiryDate:        eDate,
			CautionaryLogo:    item.CautionaryLogo,
			CautionaryText:    item.CautionaryText,
			AntidoteStatement: item.AntidoteStatement,
		}

		_, err := f.CreateFertilizer(fert)
		if err != nil {
			fmt.Printf("Failed to insert fertilizer %s: %v\n", item.ProductName, err)
		}
	}

	fmt.Println("✅ Fertilizer data loaded successfully")
	return nil
}
