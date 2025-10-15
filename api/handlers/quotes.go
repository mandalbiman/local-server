package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"bats.com/local-server/api/models"

	"github.com/gofiber/fiber/v2"
)

var quotes = []models.Quote{}

var backgroundImages = []string{
	"https://picsum.photos/seed/1/800/600",
	"https://picsum.photos/seed/2/800/600",
	"https://picsum.photos/seed/3/800/600",
	"https://picsum.photos/seed/4/800/600",
	"https://picsum.photos/seed/5/800/600",
}

func SetUpQuotesRoutes(v1 fiber.Router) {
	quotesFromFile, err := loadQuotes("sample/quotes/quotes.json")
	if err != nil {
		log.Fatal(err)
	}
	quotes = append(quotes, quotesFromFile...)
	quoteGrp := v1.Group("/quotes")
	quoteGrp.Get("/:id", getQuoteById)
	quoteGrp.Get("/", getQuotes)
}

func getQuotes(c *fiber.Ctx) error {
	return c.JSON(quotes)
}

func getQuoteById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	var selected *models.Quote
	for _, q := range quotes {
		if fmt.Sprintf("%d", q.ID) == idParam {
			selected = &q
			break
		}
	}

	if selected == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Quote not found",
		})
	}

	// Random background
	bg := backgroundImages[rand.Intn(len(backgroundImages))]
	fmt.Println(selected)

	return c.JSON(fiber.Map{
		"id":      selected.ID,
		"quote":   selected.Quote,
		"author":  selected.Author,
		"bgImage": bg,
	})
}

func loadQuotes(filename string) ([]models.Quote, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var quotes []models.Quote
	if err := json.Unmarshal(data, &quotes); err != nil {
		return nil, err
	}

	return quotes, nil
}
