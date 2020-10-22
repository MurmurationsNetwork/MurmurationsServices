package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/app"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/events/listeners"
)

type Node struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func main() {
	go listeners.NodeCreated.Listen()
	app.StartApplication()
}
