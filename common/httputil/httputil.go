package httputil

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

func GetByte(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
