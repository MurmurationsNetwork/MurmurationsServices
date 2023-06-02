package rest

import (
	geoip2 "github.com/oschwald/geoip2-golang"
)

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

type cityVO struct {
	Country string  `json:"country,omitempty"`
	City    string  `json:"city,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty"`
}

func (handler *gepIPHandler) toCityVO(record *geoip2.City) interface{} {
	res := cityVO{
		City:    record.City.Names["en"],
		Country: record.Country.IsoCode,
		Lat:     record.Location.Latitude,
		Lon:     record.Location.Longitude,
	}
	return respond{Data: res}
}
