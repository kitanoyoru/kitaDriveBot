package utils

import (
	"errors"
	"regexp"
)

func validateToken(token string) error {
	if token == "" {
		err := errors.New("Telegram token is not set!")
		return err
	}

	match, err := regexp.MatchString(`^[0-9]+:.*$`, token)
	if err != nil {
		return err
	}
	if !match {
		err := errors.New("Telegram token is incorrect. Please provide a correct Telegram Bot Token.")
		return err
	}

	return nil
}
