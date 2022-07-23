package http

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/repository/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfilesHandler interface {
	Get(c *gin.Context)
}

type profilesHandler struct {
	profileRepository db.ProfileRepository
}

func NewProfilesHandler(profileRepository db.ProfileRepository) ProfilesHandler {
	return &profilesHandler{
		profileRepository: profileRepository,
	}
}

func (handler *profilesHandler) Get(c *gin.Context) {
	profileId := c.Param("profileId")
	profile, err := handler.profileRepository.GetProfile(profileId)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}

	// remove id, __v, cuid, oid, node_id
	delete(profile, "_id")
	delete(profile, "__v")
	delete(profile, "cuid")
	delete(profile, "oid")
	delete(profile, "node_id")

	c.JSON(http.StatusOK, profile)
}
