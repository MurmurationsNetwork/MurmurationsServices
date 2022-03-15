package httputil

import (
	"fmt"
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

func Get(url string) (*http.Response, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	// IMPORTANT: defer in the sub-function to avoid err http: read on closed response body
	// defer resp.Body.Close()
	return resp, err
}

func GetByte(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []byte{}, fmt.Errorf("error the requested URL %s returned 404 not found", url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func IsValidURL(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}
