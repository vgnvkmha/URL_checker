package entities

import "google.golang.org/protobuf/types/known/timestamppb"

type Targets struct {
	ID          uint64                `json:"id"` //TODO: сделать праймари ключом
	URL         string                `json:"url"`
	IntervalSec int                   `json:"interval_sec"`
	TimeoutMS   int                   `json:"timeous_ms"`
	Active      bool                  `json:"active"`
	CreatedAt   timestamppb.Timestamp `json:"created_at"` //TODO: Сделать запись от now(), сделать параметр необязательным в JSON
	UpdatedAt   timestamppb.Timestamp `json:"updated_at"`
}
