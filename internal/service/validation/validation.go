package validation

import "errors"

func Validation(url string, interval, timeout int) error {
	return nil //TODO: сделать возврат конкретной ошибки из трёх

	// return validURL(url) &&
	// 	IsValidInterval(interval) &&
	// 	IsValidTimeout(timeout)

}

func IsValidInterval(interval int) error {
	if interval >= 10 && interval <= 3600 {
		return nil
	}
	return errors.New("Invalid interval, must be in [10;3600]")
}

func IsValidTimeout(timeout int) error {
	if timeout >= 200 && timeout <= 10_000 {
		return nil
	}
	return errors.New("Invalid timeout, must be in [200;10000]")
}

func validURL(url string) error {
	// TODO: реальная валидация URL, несколько кейсов, когда URL не валидный
	if len(url) > 0 {
		return nil
	}
	return nil
}

func ValidID(id uint64) error {
	if id >= 0 {
		return nil
	}
	return errors.New("Invalid ID")
}
