package entities

type URL struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	IntervalSec int    `json:"interval_sec"`
	TimeoutMS   int    `json:"timeous_ms"`
	Active      bool   `json:"active"`
}
