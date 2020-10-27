package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

}

func wrongSchemaFormat() {
	url := "index.murmurations.network/nodes"
	method := "POST"

	payload := strings.NewReader("{\n    \"profileUrl\": \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/exp/concurrency-issues/document.json\",\n    \"linkedSchemas\": [\n      \"455521A0657FB351689770EF3F51240C404A32F8B4026A42F056CB6A18533249\"\n    ]\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func correctSchema() {
	url := "index.murmurations.network/nodes"
	method := "POST"

	payload := strings.NewReader("{\n    \"profileUrl\": \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/exp/concurrency-issues/document.json\",\n    \"linkedSchemas\": [\n      \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/exp/concurrency-issues/schemas/correct.json\"\n    ]\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func wrongSchema() {
	url := "index.murmurations.network/nodes"
	method := "POST"

	payload := strings.NewReader("{\n    \"profileUrl\": \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/exp/concurrency-issues/document.json\",\n    \"linkedSchemas\": [\n      \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/exp/concurrency-issues/schemas/wrong.json\"\n    ]\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
