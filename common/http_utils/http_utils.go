package http_utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

var (
	client http.Client
)

func init() {
	client = http.Client{
		Timeout: 10 * time.Second,
	}
}

func GetStr(url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
