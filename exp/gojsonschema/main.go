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

	for _, schemaURL := range schemaURLs {
		schemaLoader := gojsonschema.NewReferenceLoader(schemaURL)
		result, err := gojsonschema.Validate(schemaLoader, document)
		if err != nil {
			panic(err.Error())
		}
		if result.Valid() {
			fmt.Printf("The document is valid \n")
		} else {
			fmt.Printf("The document is not valid. see errors :\n")
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
		}
	}
}
