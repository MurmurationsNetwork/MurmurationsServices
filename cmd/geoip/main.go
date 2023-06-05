package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/pkg/geoip"
)

func main() {
	logger.Info("GeoIP service starting")

	s := geoip.NewService()

	go func() {
		<-s.WaitUntilUp()
		logger.Info("GeoIP service started")
	}()

	s.Run()
}
