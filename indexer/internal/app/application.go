package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
	server = getServer()
)

func StartApplication() {
	mapUrls()
	go listen(server)

	waitForShutdown()
	cleanup()
	log.Println("Server exiting successfully")
}

func getServer() *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	return srv
}

func listen(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
