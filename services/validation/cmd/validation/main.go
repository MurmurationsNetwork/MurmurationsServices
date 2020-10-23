package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/app"
)

type Node struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func main() {
	app.StartApplication()
}
