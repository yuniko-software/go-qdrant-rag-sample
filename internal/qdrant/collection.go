package qdrant

import (
	"encoding/json"
	"fmt"
	"go-qdrant-rag-sample/internal/config"

	"github.com/go-resty/resty/v2"
)

func IsCollectionEmpty() (bool, error) {
	client := resty.New()
	host := config.QdrantHost()

	url := fmt.Sprintf("%s/collections/products/points/count", host)

	body := map[string]interface{}{}
	resp, err := client.R().
		SetQueryParam("exact", "true").
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		return false, fmt.Errorf("failed to check collection count: %w", err)
	}

	var result struct {
		Result struct {
			Count int `json:"count"`
		} `json:"result"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return false, fmt.Errorf("failed to parse count response: %w", err)
	}

	return result.Result.Count == 0, nil
}

func CreateCollection() error {
	client := resty.New()
	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     1536,
			"distance": "Cosine",
		},
	}

	host := config.QdrantHost()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(host + "/collections/products")

	if err != nil {
		return err
	}

	fmt.Println("Collection creation response:", resp.Status())
	return nil
}
