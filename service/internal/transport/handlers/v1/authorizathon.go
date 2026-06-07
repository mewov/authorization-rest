package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mewov/authorization-rest/models"
	"github.com/mewov/authorization-rest/pkg/tokens"
)

type (
	RequestRegister struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Client   string `json:"client,omitempty"`
		Role     string `json:"role,omitempty"`
	}
	RequestLogin struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
)

func HandleRegister(context *gin.Context) {
	var req RequestRegister
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error binding request",
		})
		return
	}

	lastId, err := authorization.CreateUser(req.Login, req.Email, req.Password, req.Client, req.Role)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	refresh := tokens.GenerateRefresh()
	expires := time.Now().Add(time.Hour * 24 * 7).Unix()
	if _, err := sessions.CreateSession(lastId, refresh, req.Client, expires); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error creating refresh session",
		})
		return
	}

	access, err := tokens.GenerateAccess(
		localConfig,
		lastId, req.Login, req.Email, req.Role,
	)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error creating access token",
		})
		return
	}

	context.JSON(http.StatusOK, models.DefaultResponse{
		Status:  "success",
		Message: "register successfully",
		Data: models.TokenData{
			AcessToken:   access,
			RefreshToken: refresh,
		},
	})
}

func HandleLogin(context *gin.Context) {
	var req RequestLogin
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error binding request",
		})
		return
	}

	user, err := authorization.SearchUser(req.Login, req.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	refresh := tokens.GenerateRefresh()
	expires := time.Now().Add(time.Hour * 24 * 7).Unix()
	if _, err := sessions.CreateSession(user.ID, refresh, user.Client, expires); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	access, err := tokens.GenerateAccess(
		localConfig,
		user.ID, user.Login, user.Email, user.Role,
	)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error creating access token",
		})
		return
	}

	context.JSON(http.StatusOK, models.DefaultResponse{
		Status:  "success",
		Message: "login successfully",
		Data: models.TokenData{
			AcessToken:   access,
			RefreshToken: refresh,
		},
	})
}
