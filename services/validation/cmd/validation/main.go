package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/app"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/queue"
)

type Node struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func main() {
	go queue.Listen()
	app.StartApplication()
}
