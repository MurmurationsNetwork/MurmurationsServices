package global

import "github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/config"

func Init() {
	config.Init()
	geoip2Init()
}
