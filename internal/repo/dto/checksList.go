package dto

import "time"

type ChecksList struct {
	TargetId uint64    `json:"target_id"`
	Limit    int       `json:"limit"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
}
