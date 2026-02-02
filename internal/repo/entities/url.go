package entities

type URL struct {
	URL         string `json:"url"`
	IntervalSec int    `json:"interval_sec"`
	TimeoutMS   int    `json:"timeous_ms"`
}
