package enricher

import "context"

type EnrichData struct {
	Age         int    `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type Enricher interface {
	Enrich(context.Context, string) (*EnrichData, error)
}
