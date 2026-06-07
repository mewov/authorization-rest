package transport

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mewov/authorization-rest/internal/services"
	"github.com/mewov/authorization-rest/internal/transport/handlers/v1"
	"github.com/mewov/authorization-rest/internal/transport/middleware"
	"github.com/mewov/authorization-rest/models"
)

var (
	authorization *services.AuthorizathionService
	sessions      *services.SessionsService
	localConfig   *models.LocalConfig
)

func SetResource(serviceAuth *services.AuthorizathionService, serviceSessions *services.SessionsService, config *models.LocalConfig) {
	localConfig = config
	authorization = serviceAuth
	sessions = serviceSessions
}

func Register(router *gin.Engine, logger *slog.Logger) {
	router.Use(middleware.NewLogger(logger))
	router.Use(middleware.NewRateLimit())

	v1.SetResource(authorization, sessions, localConfig)
	router.GET("/v1/status", v1.HandleStatus)

	router.POST("/v1/auth/register", v1.HandleRegister)
	router.POST("/v1/auth/login", v1.HandleLogin)
	router.POST("/v1/auth/info", v1.HandleInfo)
	router.POST("/v1/auth/logout", v1.HandleLogout)
	router.POST("/v1/auth/refresh", v1.HandleRefresh)

	fmt.Println("[+] server.register: success")
}

func Listen(router *gin.Engine, config *models.LocalConfig) {
	server := http.Server{
		Handler:      router,
		Addr:         config.SERVER_ADDRESS,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	fmt.Println("[+] server.listen: " + config.SERVER_ADDRESS + "...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("[-] server.listen:", err.Error())
	}
}
