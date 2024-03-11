package global

import (
	"fmt"
	"os"

	geoip2 "github.com/oschwald/geoip2-golang"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/config"
)

var DB *geoip2.Reader

func geoip2Init() {
	var err error
	DB, err = geoip2.Open(config.Conf.Server.DBLocation)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error when trying to Open GeoLite2-City.mmdb at %s",
				config.Conf.Server.DBLocation,
			),
			err,
		)
		os.Exit(1)
	}
}
