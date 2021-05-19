package global

import (
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/config"
	"github.com/oschwald/geoip2-golang"
)

var DB *geoip2.Reader

func geoip2Init() {
	var err error
	DB, err = geoip2.Open(config.Conf.Server.DBLocation)
	if err != nil {
		logger.Panic(fmt.Sprintf("Error when trying to Open GeoLite2-City.mmdb at %s", config.Conf.Server.DBLocation), err)
	}
}