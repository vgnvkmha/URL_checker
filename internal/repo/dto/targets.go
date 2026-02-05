package entities

import (
	"time"
)

type Targets struct {
	ID          uint64    `json:"id"` //TODO: сделать праймари ключом
	URL         string    `json:"url"`
	IntervalSec int       `json:"interval_sec"`
	TimeoutMS   int       `json:"timeous_ms"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
