package validation

func Validation(url string, interval, timeout int) bool {
	if url == "" || interval == 0 || timeout == 0 {
		return false
	}

	return validURL(url) &&
		IsValidInterval(interval) &&
		IsValidTimeout(timeout)
}

func IsValidInterval(interval int) bool {
	return interval > 10 && interval < 3600
}

func IsValidTimeout(timeout int) bool {
	return timeout > 200 && timeout < 10_000
}

func validURL(url string) bool {
	// TODO: реальная валидация URL
	return true
}

func ValidID(id int, length int) bool {
	return id < length && id >= 0
}
