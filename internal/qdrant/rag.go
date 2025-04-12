package qdrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type SearchResult struct {
	ID      interface{}    `json:"id"`
	Payload map[string]any `json:"payload"`
}

type RAGResponse struct {
	Question string         `json:"question"`
	Total    int            `json:"total"`
	Answer   []SearchResult `json:"answer"`
}

func RunRAG(question string, topK int) (RAGResponse, error) {
	results, err := SearchProducts(question, topK, nil)
	if err != nil {
		return RAGResponse{}, fmt.Errorf("retrieval error: %w", err)
	}

	// Build context from Qdrant results
	context := ""
	for _, r := range results {
		name := r.Payload["name"]
		desc := r.Payload["description"]
		context += fmt.Sprintf("- %v: %v\n", name, desc)
	}

	prompt := fmt.Sprintf(`
	You are a helpful assistant. Given the product context below, respond with a valid JSON array of the best-matching product payloads.

	Each payload must include:
	- "name": string
	- "description": string
	- "minimum_order": integer
	- "price": number
	- "price_currency": string
	- "supply_ability": integer

	DO NOT include any "id" or "score" fields.
	DO NOT wrap the response in triple backticks or Markdown formatting.

	Context:
	%s

	Question: %s

	Respond ONLY with the array of payloads.`, context, question)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return RAGResponse{}, fmt.Errorf("missing OPENAI_API_KEY")
	}

	reqBody := map[string]interface{}{
		"model": "gpt-4o-2024-08-06",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	encoded, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RAGResponse{}, fmt.Errorf("OpenAI request failed: %w", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	var raw map[string]interface{}
	if err := json.Unmarshal(resBody, &raw); err != nil {
		return RAGResponse{}, fmt.Errorf("OpenAI JSON error: %w", err)
	}

	if errObj, exists := raw["error"]; exists {
		errDetails, _ := json.MarshalIndent(errObj, "", "  ")
		return RAGResponse{}, fmt.Errorf("OpenAI API error: %s", errDetails)
	}

	choicesRaw, ok := raw["choices"].([]interface{})
	if !ok || len(choicesRaw) == 0 {
		return RAGResponse{}, fmt.Errorf("unexpected OpenAI response: %s", resBody)
	}
	message, ok := choicesRaw[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return RAGResponse{}, fmt.Errorf("unexpected format in OpenAI message field")
	}

	cleaned := strings.TrimSpace(message)
	if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
		cleaned = strings.TrimSpace(cleaned)
	}

	var payloads []map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &payloads); err != nil {
		return RAGResponse{}, fmt.Errorf("failed to parse LLM JSON: %w\nRaw:\n%s", err, cleaned)
	}

	var finalResults []SearchResult
	matchedIDs := map[interface{}]bool{}

	for _, llmPayload := range payloads {
		name, _ := llmPayload["name"].(string)
		desc, _ := llmPayload["description"].(string)

		for _, r := range results {
			rName, okName := r.Payload["name"].(string)
			rDesc, okDesc := r.Payload["description"].(string)

			if okName && okDesc && rName == name && rDesc == desc {
				if !matchedIDs[r.ID] {
					cleanPayload := map[string]any{
						"name":           rName,
						"description":    rDesc,
						"minimum_order":  r.Payload["minimum_order"],
						"price":          r.Payload["price"],
						"price_currency": r.Payload["price_currency"],
						"supply_ability": r.Payload["supply_ability"],
					}
					finalResults = append(finalResults, SearchResult{
						ID:      r.ID,
						Payload: cleanPayload,
					})
					matchedIDs[r.ID] = true
				}
				break
			}
		}
	}

	return RAGResponse{
		Question: question,
		Total:    len(finalResults),
		Answer:   finalResults,
	}, nil
}
