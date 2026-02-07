package validation

import (
	"errors"
	"net/url"
	"strings"
)

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

func validURL(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("url is empty")
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return errors.New("invalid url format")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("url must start with http or https")
	}

	if u.Host == "" {
		return errors.New("url has no host")
	}

	return nil
}
