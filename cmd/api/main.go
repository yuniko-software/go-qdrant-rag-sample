package main

import (
	"fmt"
	"log"
	"os"

	"go-qdrant-rag-sample/internal/api"
	"go-qdrant-rag-sample/internal/config"
	"go-qdrant-rag-sample/internal/qdrant"
	"go-qdrant-rag-sample/internal/utils"
)

func main() {

	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = "./env/.env"
	}

	fmt.Println("Loading env file:", envFile)

	if err := config.LoadEnv(envFile); err != nil {
		log.Fatalf("Failed to load env file %s: %v", envFile, err)
	}

	if err := qdrant.CreateCollection(); err != nil {
		log.Fatalf("Failed to create Qdrant collection: %v", err)
	}

	isEmpty, err := qdrant.IsCollectionEmpty()
	if err != nil {
		log.Fatalf("Failed to check if collection is empty: %v", err)
	}

	if isEmpty {
		fmt.Println("Collection is empty. Loading products and inserting...")

		// Now load products ONLY IF we need them
		products, err := utils.LoadProductsCSV("data/products.csv")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Loaded %d products\n", len(products))
		qdrant.InsertAllProducts(products)
	} else {
		fmt.Println("Collection already has data. Skipping insert.")
	}

	api.SetupAPI()
}
