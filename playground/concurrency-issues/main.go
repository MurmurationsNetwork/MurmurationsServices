package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	wrongSchemaFormat()
	correctSchema()
	wrongSchema()
}

func wrongSchemaFormat() {
	url := "http://index.murmurations.network/v1/nodes"
	method := "POST"

	payload := strings.NewReader("{\n    \"profileUrl\": \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/playground/concurrency-issues/document.json\",\n    \"linkedSchemas\": [\n      \"455521A0657FB351689770EF3F51240C404A32F8B4026A42F056CB6A18533249\"\n    ]\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("==================================")
		fmt.Printf("err %+v \n", err)
		fmt.Println("==================================")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func correctSchema() {
	url := "http://index.murmurations.network/v1/nodes"
	method := "POST"

	payload := strings.NewReader("{\n    \"profileUrl\": \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/playground/concurrency-issues/document.json\",\n    \"linkedSchemas\": [\n      \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/playground/concurrency-issues/schemas/correct.json\"\n    ]\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("==================================")
		fmt.Printf("err %+v \n", err)
		fmt.Println("==================================")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func wrongSchema() {
	url := "http://index.murmurations.network/v1/nodes"
	method := "POST"

	payload := strings.NewReader("{\n    \"profileUrl\": \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/playground/concurrency-issues/document.json\",\n    \"linkedSchemas\": [\n      \"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/indexAPI/playground/concurrency-issues/schemas/wrong.json\"\n    ]\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("==================================")
		fmt.Printf("err %+v \n", err)
		fmt.Println("==================================")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
