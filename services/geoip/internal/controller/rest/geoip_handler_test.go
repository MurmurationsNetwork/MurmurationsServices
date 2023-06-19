package rest_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/controller/rest"
)

type fakeGeoIPService struct{}

func (f fakeGeoIPService) GetCity(_ string) (*geoip2.City, resterr.RestErr) {
	names := map[string]string{
		"en": "Taiwan",
	}
	return &geoip2.City{
		City: struct {
			Names     map[string]string "maxminddb:\"names\""
			GeoNameID uint              "maxminddb:\"geoname_id\""
		}{
			Names: names,
		},
	}, nil
}
func TestGetCity(t *testing.T) {
	tests := map[string]struct {
		param   gin.Params
		bodyStr string
	}{
		"success": {
			param: []gin.Param{
				{
					Key:   "ip",
					Value: "1.164.203.137",
				},
			},
			bodyStr: `{"data":{"city":"Taiwan"}}`,
		},
		"wrong param": {
			param: []gin.Param{
				{
					Key:   "wrong-key",
					Value: "1.164.203.137",
				},
			},
			bodyStr: `{"message":"Invalid IP address.","status":400}`,
		},
	}

	for _, tc := range tests {
		geoIPService := fakeGeoIPService{}
		geoIPHandler := rest.NewGeoIPHandler(geoIPService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = tc.param

		geoIPHandler.GetCity(c)

		assert.Equal(t, tc.bodyStr, w.Body.String())
	}
}
