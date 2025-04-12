package qdrant

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func GetEmbedding(text string) ([]float32, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing OPENAI_API_KEY")
	}

	client := resty.New()
	requestBody := map[string]interface{}{
		"input": text,
		"model": "text-embedding-ada-002",
	}

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post("https://api.openai.com/v1/embeddings")

	if err != nil {
		return nil, err
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("embedding response was empty")
	}

	return result.Data[0].Embedding, nil
}
