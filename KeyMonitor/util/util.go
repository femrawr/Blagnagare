package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"bytes"
)

var url string

func RegisterUrl(toRegister string) {
	url = toRegister
}

func PostToUrl(data string) error {
	if len(data) <= 0 {
		return errors.New("data is empty")
	}

	if len(url) <= 0 {
		return errors.New("url is empty")
	}

	payload := map[string] string {
		"content": data,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to send")
	}

	return nil
}