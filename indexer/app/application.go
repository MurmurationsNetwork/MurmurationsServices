package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/indexer/datasources/mongo/nodes_db"
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

func waitForShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

func cleanup() {
	nodes_db.Disconnect()
}
