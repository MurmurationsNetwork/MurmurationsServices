package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func GetWithBearerToken(url string, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func GetByteWithBearerToken(url string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []byte{}, fmt.Errorf("error the requested URL %s returned 404 not found", url)
	}

	// Read the response body into a byte array
	data, err := io.ReadAll(resp.Body)
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

func GetJSONStr(source string) (string, error) {
	jsonByte, err := GetByte(source)
	if err != nil {
		return "", err
	}
	buffer := bytes.Buffer{}
	err = json.Compact(&buffer, jsonByte)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
