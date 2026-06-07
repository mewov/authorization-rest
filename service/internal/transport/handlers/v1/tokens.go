package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mewov/authorization-rest/models"
	"github.com/mewov/authorization-rest/pkg/tokens"
)

type (
	RequestToken struct {
		Token string `json:"token"`
	}
)

func HandleRefresh(context *gin.Context) {
	var req RequestToken
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error binding request",
		})
		return
	}

	session, err := sessions.SearchSession(req.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	user, err := authorization.SearchUserByID(session.UserID)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "account is invalid, id is not found: " + fmt.Sprint(session.UserID),
		})
		return
	}

	err = sessions.RemoveSession(req.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error remove session: " + fmt.Sprint(session.UserID),
		})
		return
	}

	refresh := tokens.GenerateRefresh()
	expires := time.Now().Add(time.Hour * 24 * 7).Unix()
	if _, err := sessions.CreateSession(user.ID, refresh, user.Client, expires); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error creating refresh session",
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

func HandleLogout(context *gin.Context) {
	var req RequestToken
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error binding request",
		})
		return
	}

	err := sessions.RemoveSession(req.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.DefaultResponse{
		Status:  "success",
		Message: "logout successfully",
	})
}

func HandleInfo(context *gin.Context) {
	var req RequestToken
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "error binding request",
		})
		return
	}

	access, err := tokens.CheckAccess(localConfig, req.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status:  "error",
			Message: "access token is invalid",
		})
		return
	}

	context.JSON(http.StatusOK, models.DefaultResponse{
		Status:  "success",
		Message: "user info successfully",
		Data: models.UserData{
			UserID: access.UserId,
			Login:  access.Login,
			Email:  access.Email,
			Role:   access.Role,
		},
	})
}
