package eel

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sync"

	"github.com/born2ngopi/eel/pkg/client"
	"github.com/born2ngopi/eel/pkg/memcache"
)

type ell struct {
	token string
	host  string
}

var e *ell
var once sync.Once

func Init(token, host string) error {

	var errRes error
	once.Do(func() {

		var body = struct {
			Token string `json:"token"`
		}{
			Token: token,
		}

		b, err := json.Marshal(body)
		if err != nil {
			errRes = err
		}

		uri, err := url.JoinPath(host, "/v1/service/validate")
		if err != nil {
			errRes = err
		}

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(b))
		if err != nil {
			errRes = err
		}

		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errRes = err
		}

		if resp.StatusCode != http.StatusOK {
			errRes = errors.New("invalid token")
		}

		e = &ell{
			token: token,
			host:  host,
		}
	})

	return errRes
}

func GetToken(key string) (string, error) {
	token := memcache.Get(key)
	if token != "" {
		return token, nil
	}

	// get from ell
	token, err := client.Do(e.token, e.host, key)
	if err != nil {
		return "", err
	}

	memcache.Set(key, token)

	return token, nil
}
