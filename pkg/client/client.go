package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func Do(token, host, key string) (string, error) {

	uri, err := url.JoinPath(host, fmt.Sprintf("/v1/service/%s/token", key))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid token")
	}

	var res = struct {
		Token string `json:"token"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Token, nil
}
