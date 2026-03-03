package mapper

import (
	entities "URL_checker/internal/repo/dto"
	"encoding/json"
)

func FromTarget(dto entities.Targets) ([]byte, error) {
	return json.Marshal(dto)
}

func ToTarget(data []byte) (entities.Targets, error) {
	var dto entities.Targets
	err := json.Unmarshal(data, &dto)
	return dto, err
}
