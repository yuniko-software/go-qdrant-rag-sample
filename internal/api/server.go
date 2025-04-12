package api

import (
	"log"
	"strconv"

	"go-qdrant-rag-sample/internal/qdrant"

	"github.com/gofiber/fiber/v2"
)

func SetupAPI() {
	app := fiber.New()

	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.Status(400).SendString("Missing query param `q`")
		}

		topK := 5
		if top := c.Query("top"); top != "" {
			if parsed, err := strconv.Atoi(top); err == nil && parsed > 0 {
				topK = parsed
			}
		}

		var maxPrice *float64
		if price := c.Query("maxPrice"); price != "" {
			if p, err := strconv.ParseFloat(price, 64); err == nil {
				maxPrice = &p
			}
		}

		results, err := qdrant.SearchProducts(query, topK, maxPrice)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(results)
	})

	app.Get("/rag", func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.Status(400).SendString("Missing query param `q`")
		}

		topK := 5
		if top := c.Query("top"); top != "" {
			if parsed, err := strconv.Atoi(top); err == nil && parsed > 0 {
				topK = parsed
			}
		}

		response, err := qdrant.RunRAG(query, topK)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(response)
	})

	log.Fatal(app.Listen(":8080"))
}
