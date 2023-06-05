package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/service"
)

type GeoIPHandler interface {
	GetCity(c *gin.Context)
}

type geoIPHandler struct {
	svc service.GeoIPService
}

func NewGeoIPHandler(svc service.GeoIPService) GeoIPHandler {
	return &geoIPHandler{
		svc: svc,
	}
}

func (handler *geoIPHandler) GetCity(c *gin.Context) {
	ip, found := c.Params.Get("ip")
	if !found {
		restErr := resterr.NewBadRequestError("Invalid IP address.")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := handler.svc.GetCity(ip)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, handler.toCityVO(result))
}
