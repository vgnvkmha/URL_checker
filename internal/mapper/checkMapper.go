package mapper

import (
	entities "URL_checker/internal/repo/dto"
	"encoding/json"
)

func ToCheck(data []byte) (entities.Checks, error) {
	var dto entities.Checks
	err := json.Unmarshal(data, &dto)
	return dto, err
}

func FromCheck(dto entities.Checks) ([]byte, error) {
	return json.Marshal(dto)
}
