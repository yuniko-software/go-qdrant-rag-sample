package qdrant

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"go-qdrant-rag-sample/internal/config"
	"go-qdrant-rag-sample/internal/models"

	"github.com/go-resty/resty/v2"
)

func InsertProduct(id string, embedding []float32, p models.Product) error {
	client := resty.New()

	body := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":     id,
				"vector": embedding,
				"payload": map[string]interface{}{
					"name":           p.Name,
					"description":    p.Description,
					"price":          p.Price,
					"price_currency": p.PriceCurrency,
					"supply_ability": p.SupplyAbility,
					"minimum_order":  p.MinimumOrder,
				},
			},
		},
	}

	host := config.QdrantHost()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(host + "/collections/products/points")

	if err != nil {
		return err
	}

	fmt.Printf("Inserted %s: %s\n", id, resp.Status())
	return nil
}

func InsertAllProducts(products []models.Product) {
	fmt.Printf("Inserting %d products with parallel processing...\n", len(products))

	workerCount := runtime.NumCPU() * 2
	jobs := make(chan models.Product, workerCount)
	wg := sync.WaitGroup{}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for p := range jobs {
				input := p.ToEmbeddingInput()
				embedding, err := GetEmbedding(input)
				if err != nil {
					log.Printf("[Worker %d] Failed to embed product %s: %v", workerID, p.ID, err)
					continue
				}
				err = InsertProduct(p.ID, embedding, p)
				if err != nil {
					log.Printf("[Worker %d] Failed to insert %s: %v", workerID, p.ID, err)
				} else {
					log.Printf("[Worker %d] Inserted product: %s", workerID, p.ID)
				}
			}
		}(i + 1)
	}

	for _, p := range products {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
	fmt.Println("All products processed.")
}
