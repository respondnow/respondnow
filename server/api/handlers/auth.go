package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/respondnow/respond/server/pkg/auth"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	auth2 "github.com/respondnow/respond/server/pkg/database/mongodb/auth"
	"github.com/respondnow/respond/server/utils"
	"github.com/sirupsen/logrus"
)

// SignUp godoc
//
//	@Summary		Signup to RespondNow
//	@Description	Signup to RespondNow
//	@id				SignUp
//
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			"signup"	body		auth.AddUserInput	true	"Signup to RespondNow"
//	@Success		200					{object}	auth.SignupResponseDTO
//	@Failure		400					{object}	utils.DefaultResponseDTO
//	@Failure		404					{object}	utils.DefaultResponseDTO
//	@Failure		500					{object}	utils.DefaultResponseDTO
//	@Router			/auth/signup [post]
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = auth.SignupResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		var user auth.AddUserInput
		if err := c.ShouldBindJSON(&user); err != nil {
			response.DefaultResponseDTO.Message = "Invalid request body"
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		logFields := logrus.Fields{
			"correlationId": response.DefaultResponseDTO.CorrelationId,
		}

		if err := validator.New().Struct(user); err != nil {
			logrus.WithFields(logFields).WithError(err).Error("failed to validate the request")
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		timeOut := time.Second * 10
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		_, err := auth.NewAuthService(auth2.NewAuthOperator(mongodb.Operator)).Signup(ctx, user)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				response.DefaultResponseDTO.Message = "Signup timed out"
			} else {
				response.DefaultResponseDTO.Message = "Signup failed: " + err.Error()
			}
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = string(utils.SUCCESS)
		response.DefaultResponseDTO.Message = "User registered successfully"
		c.JSON(http.StatusOK, response)
	}
}

// Login godoc
//
//	@Summary		Login to RespondNow
//	@Description	Login to RespondNow
//	@id				Login
//
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			"login"	body		auth.LoginUserInput	true	"Login to RespondNow"
//	@Success		200					{object}	auth.LoginResponseDTO
//	@Failure		400					{object}	utils.DefaultResponseDTO
//	@Failure		404					{object}	utils.DefaultResponseDTO
//	@Failure		500					{object}	utils.DefaultResponseDTO
//	@Router			/auth/login [post]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = auth.LoginResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		var loginReq auth.LoginUserInput
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			response.DefaultResponseDTO.Message = "invalid request body"
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		logFields := logrus.Fields{
			"correlationId": response.DefaultResponseDTO.CorrelationId,
		}

		if err := validator.New().Struct(loginReq); err != nil {
			logrus.WithFields(logFields).WithError(err).Error("failed to validate the request")
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		timeOut := time.Second * 10
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		authService := auth.NewAuthService(auth2.NewAuthOperator(mongodb.Operator))
		user, err := authService.Login(ctx, loginReq)
		if err != nil {
			response.DefaultResponseDTO.Message = "Login failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		if user.ChangePasswordRequired {
			response.DefaultResponseDTO.Message = "Change Password is required"
			response.Data.ChangeUserPassword = user.ChangePasswordRequired
			response.Status = string(utils.SUCCESS)
			c.JSON(http.StatusOK, response)
			return
		}

		token, err := authService.CreateJWTToken(user.Email, user.UserID, user.Name)
		if err != nil {
			response.Status = string(utils.ERROR)
			response.DefaultResponseDTO.Message = "failed to generate token: " + err.Error()
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response.Status = string(utils.SUCCESS)
		response.DefaultResponseDTO.Message = "Login successful"
		response.Data.Token = token
		c.JSON(http.StatusOK, response)
	}
}

// ChangePassword godoc
//
//	@Summary		ChangePassword of RespondNow
//	@Description	ChangePassword of RespondNow
//	@id				ChangePassword
//
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			"changePassword"	body		auth.ChangeUserPasswordInput	true	"ChangePassword of RespondNow"
//	@Success		200					{object}	auth.ChangePasswordResponseDTO
//	@Failure		400					{object}	utils.DefaultResponseDTO
//	@Failure		404					{object}	utils.DefaultResponseDTO
//	@Failure		500					{object}	utils.DefaultResponseDTO
//	@Router			/auth/changePassword [post]
func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = auth.ChangePasswordResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		var changePasswordReq auth.ChangeUserPasswordInput
		if err := c.ShouldBindJSON(&changePasswordReq); err != nil {
			response.DefaultResponseDTO.Message = "invalid request body"
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		logFields := logrus.Fields{
			"correlationId": response.DefaultResponseDTO.CorrelationId,
		}

		if err := validator.New().Struct(changePasswordReq); err != nil {
			logrus.WithFields(logFields).WithError(err).Error("failed to validate the request")
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		timeOut := time.Second * 10
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		authService := auth.NewAuthService(auth2.NewAuthOperator(mongodb.Operator))
		err := authService.ChangePassword(ctx, changePasswordReq)
		if err != nil {
			response.DefaultResponseDTO.Message = "Login failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		user, err := authService.Login(ctx, auth.LoginUserInput{
			Email:    changePasswordReq.Email,
			Password: changePasswordReq.NewPassword,
		})
		if err != nil {
			response.DefaultResponseDTO.Message = "Login failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		token, err := authService.CreateJWTToken(user.Email, user.UserID, user.Name)
		if err != nil {
			response.DefaultResponseDTO.Message = "failed to generate token: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response.DefaultResponseDTO.Message = "Password has been changed"
		response.Data.Token = token
		response.Status = string(utils.SUCCESS)
		c.JSON(http.StatusOK, response)
	}
}
