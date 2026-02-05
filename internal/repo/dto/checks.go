package dto

import (
	"time"
)

type Checks struct {
	ID         uint64    `json:"id"`
	TargetId   uint64    `json:"target_id"`
	CheckedAt  time.Time `json:"checked_at"`
	OK         bool      `json:"ok"`
	StatusCode uint8     `json:"status_code"`
	LatencyMs  uint8     `json:"latency_ms"`
	Error      string    `json:"error"`
}
