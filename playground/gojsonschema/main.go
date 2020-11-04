package main

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

var (
	documentURL = "https://ic3.dev/test1.json"
	schemaURLs  = []string{"https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsLibrary/master/schemas/demo-v1.json"}
)

func main() {
	document := gojsonschema.NewReferenceLoader(documentURL)
	data, _ := document.LoadJSON()
	linkedSchemas := data.(map[string]interface{})["linkedSchemas"].([]interface{})
	for _, schema := range linkedSchemas {
		fmt.Printf("%+v \n", schema.(string))
	}
}
