package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/controllers/http/ping"

func mapUrls() {
	router.GET("/ping", ping.Ping)
}
