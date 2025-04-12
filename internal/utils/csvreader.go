package utils

import (
	"encoding/csv"
	"os"
	"go-qdrant-rag-sample/internal/models"
	"strconv"
)

func LoadProductsCSV(filename string) ([]models.Product, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	reader.FieldsPerRecord = -1

	raws, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var products []models.Product
	for i, row := range raws {
		if i == 0 || len(row) < 7 {
			continue
		}

		price, _ := strconv.ParseFloat(row[3], 64)
		supply, _ := strconv.Atoi(row[5])
		minOrder, _ := strconv.Atoi(row[6])

		products = append(products, models.Product{
			ID:            row[0],
			Name:          row[1],
			Description:   row[2],
			Price:         price,
			PriceCurrency: row[4],
			SupplyAbility: supply,
			MinimumOrder:  minOrder,
		})
	}

	return products, nil
}
