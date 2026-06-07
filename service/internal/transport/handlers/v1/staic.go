package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mewov/authorization-rest/internal/services"
	"github.com/mewov/authorization-rest/models"
)

var (
	authorization *services.AuthorizathionService
	sessions      *services.SessionsService
	localConfig   *models.LocalConfig
)

func SetResource(serviceAuth *services.AuthorizathionService, serviceSessions *services.SessionsService, config *models.LocalConfig) {
	authorization = serviceAuth
	sessions = serviceSessions
	localConfig = config
}

func HandleStatus(ctx *gin.Context) {
	statusAuth := "is not work"
	if authorization != nil {
		statusAuth = "is work"
	}

	statusSession := "is not work"
	if sessions != nil {
		statusSession = "is work"
	}

	ctx.JSON(200, models.DefaultResponse{
		Status:  "success",
		Message: "status services",
		Data: models.StatusData{
			StatusAuth:    statusAuth,
			StatusSession: statusSession,
		},
	})
}
