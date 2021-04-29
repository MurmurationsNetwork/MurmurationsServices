package http

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/service"
	"github.com/gin-gonic/gin"
)

type GepIPHandler interface {
	GetCity(c *gin.Context)
}

type gepIPHandler struct {
	svc service.GeoIPService
}

func NewGepIPHandler(svc service.GeoIPService) GepIPHandler {
	return &gepIPHandler{
		svc: svc,
	}
}

func (handler *gepIPHandler) GetCity(c *gin.Context) {
	ip, found := c.Params.Get("ip")
	if !found {
		restErr := resterr.NewBadRequestError("Invalid ip.")
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
