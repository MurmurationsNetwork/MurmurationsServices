package http

import (
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/db"
	"github.com/gin-gonic/gin"
)

type UpdatesHandler interface {
	Get(c *gin.Context)
}

type updatesHandler struct {
	updateRepository db.UpdateRepository
}

func NewUpdatesHandler(updateRepository db.UpdateRepository) UpdatesHandler {
	return &updatesHandler{
		updateRepository: updateRepository,
	}
}

func (handler *updatesHandler) Get(c *gin.Context) {
	schemaName := c.Param("schemaName")
	update, err := handler.updateRepository.GetUpdate(schemaName)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, update)
}
