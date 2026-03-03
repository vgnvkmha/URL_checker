package mapper

// import (
// 	"encoding/base64"
// 	"encoding/json"
// 	"strings"
// )

// func DecodeValue(val string) (string, error) {
// 	// убираем возможные лишние кавычки
// 	clean := strings.Trim(val, "\"")

// 	// декодируем base64
// 	decoded, err := base64.StdEncoding.DecodeString(clean)
// 	if err != nil {
// 		return "", err
// 	}

// 	// парсим JSON
// 	var result string
// 	if err := json.Unmarshal(decoded, &result); err != nil {
// 		return "", err
// 	}

// 	return result, nil
// }
