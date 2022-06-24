package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/geoip/internal/controller/http"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"
)

type fakeGeoIPService struct{}

func (f fakeGeoIPService) GetCity(ip string) (*geoip2.City, resterr.RestErr) {
	names := map[string]string{
		"en": "Taiwan",
	}
	return &geoip2.City{
		City: struct {
			GeoNameID uint              "maxminddb:\"geoname_id\""
			Names     map[string]string "maxminddb:\"names\""
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
			bodyStr: `{"message":"Invalid ip.","status":400}`,
		},
	}

	for _, tc := range tests {
		geoIPService := fakeGeoIPService{}
		gepIPHandler := http.NewGepIPHandler(geoIPService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = tc.param

		gepIPHandler.GetCity(c)

		assert.Equal(t, tc.bodyStr, w.Body.String())
	}
}
