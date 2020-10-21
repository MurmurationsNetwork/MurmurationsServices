package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/controllers/ping"

func mapUrls() {
	router.GET("/ping", ping.Ping)
}
