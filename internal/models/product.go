package models

import "fmt"

type Product struct {
	ID            string
	Name          string
	Description   string
	Price         float64
	PriceCurrency string
	SupplyAbility int
	MinimumOrder  int
}

func (p Product) ToEmbeddingInput() string {
	return fmt.Sprintf(
		"%s. %s. The price is %.2f %s. Minimum order: %d units. Supply ability: %d units.",
		p.Name,
		p.Description,
		p.Price,
		p.PriceCurrency,
		p.MinimumOrder,
		p.SupplyAbility,
	)
}
