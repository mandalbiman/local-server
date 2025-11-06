package utils

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

func AlphaSpaceDot(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[A-Za-z .]+$`)
	return re.MatchString(fl.Field().String())
}

func IsValidDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String()) // Expected format: YYYY-MM-DD
	return err == nil
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
