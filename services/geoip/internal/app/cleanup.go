package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/global"

func cleanup() {
	global.DB.Close()
}
