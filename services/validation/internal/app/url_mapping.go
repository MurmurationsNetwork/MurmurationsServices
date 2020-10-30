package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/controller/http"

func mapUrls() {
	router.GET("/ping", http.PingController.Ping)
}
