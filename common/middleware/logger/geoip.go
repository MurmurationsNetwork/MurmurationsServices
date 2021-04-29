package logger

import (
	"encoding/json"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/gin-gonic/gin"
)

type respond struct {
	Data geoInfo `json:"data,omitempty"`
}

type geoInfo struct {
	City    string  `json:"city,omitempty"`
	Country string  `json:"country,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty"`
}

func getGeoInfo(c *gin.Context) *geoInfo {
	bytes, err := httputil.GetByte(fmt.Sprintf("http://geoip-app:8080/city/%s", c.ClientIP()))
	if err != nil {
		return &geoInfo{}
	}

	var respondData respond

	err = json.Unmarshal(bytes, &respondData)
	if err != nil {
		return &geoInfo{}
	}

	return &respondData.Data
}
