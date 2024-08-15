package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/respondnow/respond/server/utils"
	"github.com/respondnow/respond/server/version"
	"github.com/sirupsen/logrus"
)

// StatusHandler godoc
//
//	@Summary		Status of RespondNow Server
//	@Description	Status of RespondNow Server
//
//	@id				Status
//
//	@Tags			Miscellaneous
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	APIStatus
//	@Failure		404	{object}	utils.DefaultResponseDTO
//	@Router			/status [get]
//
// StatusHandler handles the status of server.
func StatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, APIStatus{Status: "up"})
	}
}

type APIStatus struct {
	Status string `json:"status"`
}

// InitVersionInfo godoc
//
//	@Summary		Version of RespondNow Server
//	@Description	Version of RespondNow Server
//
//	@id				Version
//
//	@Tags			Miscellaneous
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	version.Version
//	@Failure		404	{object}	utils.DefaultResponseDTO
//	@Failure		500	{object}	utils.DefaultResponseDTO
//	@Router			/version [get]
//
// InitVersionInfo handles the version of server.
func InitVersionInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := utils.NewUtils().RandStringBytes(10)
		versionConfigs, err := version.GetVersionInfo()
		if err != nil {
			logrus.WithField("correlationId", correlationId).Error(err)
			errMsg := utils.DefaultResponseDTO{
				Status:        string(utils.ERROR),
				Message:       err.Error(),
				CorrelationId: correlationId,
			}
			c.JSON(http.StatusInternalServerError, errMsg)
		}

		c.JSON(http.StatusOK, versionConfigs)
	}
}
