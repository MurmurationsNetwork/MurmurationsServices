package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"

	"github.com/olivere/elastic/v7"
)

func Init() {
	config.Init()
}

func main() {
	Init()

	count := 0
	updateSize := 100

	fmt.Println("ES URL", config.Conf.ES.URL)
	client, err := elastic.NewClient(elastic.SetURL(config.Conf.ES.URL))
	if err != nil {
		fmt.Println("Error when trying to ping Elasticsearch", err)
		return
	}

	ctx := context.Background()

	scrollService := client.Scroll("nodes").
		Size(updateSize).
		Query(elastic.NewMatchAllQuery())
	searchResult, err := scrollService.Do(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Process the initial set of documents
	iterateUpdate(ctx, client, searchResult, &count)

	// Continue scrolling and updating
	for {
		scrollService := client.Scroll("nodes").
			Size(updateSize).
			ScrollId(searchResult.ScrollId)
		searchResult, err = scrollService.Do(ctx)
		if err == io.EOF {
			// No more documents to process
			break
		}
		if err != nil {
			fmt.Println("Scroll error")
			fmt.Println("current scroll id", searchResult.ScrollId)
			fmt.Println(err)
			return
		}

		if searchResult.Hits == nil || len(searchResult.Hits.Hits) == 0 {
			// No more documents to process
			break
		}

		// Process the next set of documents
		iterateUpdate(ctx, client, searchResult, &count)
	}

	// Clear the scroll context
	clearScrollService := client.ClearScroll(searchResult.ScrollId)
	_, err = clearScrollService.Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Update complete with total", count, "nodes")
}

func iterateUpdate(
	ctx context.Context,
	client *elastic.Client,
	searchResult *elastic.SearchResult,
	count *int,
) {
	for _, hit := range searchResult.Hits.Hits {
		bytes, err := hit.Source.MarshalJSON()
		if err != nil {
			fmt.Println(err)
			continue
		}
		var result map[string]interface{}
		if err := json.Unmarshal(bytes, &result); err != nil {
			fmt.Println(err)
			continue
		}

		profileURL, ok := result["profile_url"].(string)
		if !ok {
			fmt.Println("profile_url is not a string or not exist")
			fmt.Println(result)
			continue
		}

		// Get data from profile_url
		name, err := getProfileName(profileURL)
		if err != nil {
			fmt.Println(
				"get profile name error with url: ",
				profileURL,
				"Error: ",
				err,
			)
			continue
		}

		update := elastic.NewUpdateService(client).
			Index(hit.Index).
			Id(hit.Id).
			Doc(map[string]interface{}{
				"name": name,
			})
		_, err = update.Do(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}

		*count++
		fmt.Println("Count:", *count, "with Id:", hit.Id, "and name:", name)
	}
}

func getProfileName(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// get name from resp.Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	name, ok := data["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is not a string or not exist")
	}

	return name, nil
}
