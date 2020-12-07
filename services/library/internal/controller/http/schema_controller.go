package http

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	SchemaController schemaControllerInterface = &schemaController{}
)

type schemaControllerInterface interface {
	Search(c *gin.Context)
}

type schemaController struct{}

func (cont *schemaController) Search(c *gin.Context) {
	searchRes, err := service.SchemaService.Search()
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, searchRes.Marshall())
}
