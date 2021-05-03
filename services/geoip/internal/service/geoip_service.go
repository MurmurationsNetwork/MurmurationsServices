package service

import (
	"errors"
	"fmt"
	"net"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/global"
	"github.com/oschwald/geoip2-golang"
)

type GeoIPService interface {
	GetCity(string) (*geoip2.City, resterr.RestErr)
}

type geoIPService struct {
}

func NewGeoIPService() GeoIPService {
	return &geoIPService{}
}

func (s *geoIPService) GetCity(ipStr string) (*geoip2.City, resterr.RestErr) {
	ip := net.ParseIP(ipStr)

	record, err := global.DB.City(ip)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to get geographic info from an IP geress: %s", ipStr), err)
		return nil, resterr.NewInternalServerError("Error when trying get geographic info.", errors.New("database error"))
	}

	return record, nil
}
