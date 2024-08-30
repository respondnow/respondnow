package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	hierarchy2 "github.com/respondnow/respondnow/server/pkg/hierarchy"

	"github.com/respondnow/respondnow/server/pkg/database/mongodb/hierarchy"
	"github.com/respondnow/respondnow/server/pkg/user"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/respondnow/respondnow/server/pkg/database/mongodb"
	auth2 "github.com/respondnow/respondnow/server/pkg/database/mongodb/user"
	"github.com/respondnow/respondnow/server/utils"
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
//	@Param			"signup"		body		user.AddUserInput	true	"Signup to RespondNow"
//	@Param			correlationId	query		string				false	"correlationId"
//	@Success		200				{object}	user.SignupResponseDTO
//	@Failure		400				{object}	utils.DefaultResponseDTO
//	@Failure		404				{object}	utils.DefaultResponseDTO
//	@Failure		500				{object}	utils.DefaultResponseDTO
//	@Router			/auth/signup [post]
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = user.SignupResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		var u user.AddUserInput
		if err := c.ShouldBindJSON(&u); err != nil {
			response.DefaultResponseDTO.Message = "Invalid request body"
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		logFields := logrus.Fields{
			"correlationId": response.DefaultResponseDTO.CorrelationId,
		}

		if err := validator.New().Struct(u); err != nil {
			logrus.WithFields(logFields).WithError(err).Error("failed to validate the request")
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		timeOut := time.Second * 10
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		_, err := user.NewAuthService(auth2.NewAuthOperator(mongodb.Operator)).Signup(ctx, u)
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
//	@Param			"login"			body		user.LoginUserInput	true	"Login to RespondNow"
//	@Param			correlationId	query		string				false	"correlationId"
//	@Success		200				{object}	user.LoginResponseDTO
//	@Failure		400				{object}	utils.DefaultResponseDTO
//	@Failure		404				{object}	utils.DefaultResponseDTO
//	@Failure		500				{object}	utils.DefaultResponseDTO
//	@Router			/auth/login [post]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = user.LoginResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		var loginReq user.LoginUserInput
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

		authService := user.NewAuthService(auth2.NewAuthOperator(mongodb.Operator))
		user, err := authService.Login(ctx, loginReq)
		if err != nil {
			response.DefaultResponseDTO.Message = "Login failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		token, err := authService.CreateJWTToken(user.Email, user.UserID, user.Name)
		if err != nil {
			response.Status = string(utils.ERROR)
			response.DefaultResponseDTO.Message = "failed to generate token: " + err.Error()
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		if user.ChangePasswordRequired {
			response.DefaultResponseDTO.Message = "Change Password is required"
			response.Data.ChangeUserPassword = user.ChangePasswordRequired
			response.Data.LastLoginAt = user.LastLoginAt
			response.Data.Token = token
			response.Status = string(utils.SUCCESS)
			c.JSON(http.StatusOK, response)
			return
		}

		err = authService.UpdateLastLogin(ctx, loginReq)
		if err != nil {
			response.DefaultResponseDTO.Message = "Login failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		response.Status = string(utils.SUCCESS)
		response.DefaultResponseDTO.Message = "Login successful"
		response.Data.Token = token
		response.Data.LastLoginAt = user.LastLoginAt
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
//	@Param			"changePassword"	body		user.ChangeUserPasswordInput	true	"ChangePassword of RespondNow"
//	@Param			correlationId		query		string							false	"correlationId"
//	@Success		200					{object}	user.ChangePasswordResponseDTO
//	@Failure		400					{object}	utils.DefaultResponseDTO
//	@Failure		404					{object}	utils.DefaultResponseDTO
//	@Failure		500					{object}	utils.DefaultResponseDTO
//	@Router			/auth/changePassword [post]
func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = user.ChangePasswordResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		var changePasswordReq user.ChangeUserPasswordInput
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

		authService := user.NewAuthService(auth2.NewAuthOperator(mongodb.Operator))
		err := authService.ChangePassword(ctx, changePasswordReq)
		if err != nil {
			response.DefaultResponseDTO.Message = "Change Password Failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		u, err := authService.Login(ctx, user.LoginUserInput{
			Email:    changePasswordReq.Email,
			Password: changePasswordReq.NewPassword,
		})
		if err != nil {
			response.DefaultResponseDTO.Message = "Change Password Failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		token, err := authService.CreateJWTToken(u.Email, u.UserID, u.Name)
		if err != nil {
			response.DefaultResponseDTO.Message = "failed to generate token: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		err = authService.UpdateLastLogin(ctx, user.LoginUserInput{
			Email:    u.Email,
			Password: changePasswordReq.NewPassword,
		})
		if err != nil {
			response.DefaultResponseDTO.Message = "Change Password failed: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		response.DefaultResponseDTO.Message = "Password has been changed"
		response.Data.Token = token
		response.Data.LastLoginAt = u.LastLoginAt
		response.Status = string(utils.SUCCESS)
		c.JSON(http.StatusOK, response)
	}
}

// GetUserMapping godoc
//
//	@Summary		GetUserMapping of RespondNow
//	@Description	GetUserMapping of RespondNow
//	@id				GetUserMapping
//
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Param			correlationId	query		string	false	"correlationId"
//	@Param			userId			query		string	false	"userId"
//	@Success		200				{object}	user.GetUserMappingResponseDTO
//	@Failure		400				{object}	utils.DefaultResponseDTO
//	@Failure		404				{object}	utils.DefaultResponseDTO
//	@Failure		500				{object}	utils.DefaultResponseDTO
//	@Router			/auth/userMapping [get]
func GetUserMapping() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = user.GetUserMappingResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		userID := c.Query("userId")
		if userID == "" {
			response.DefaultResponseDTO.Message = "userId is required in the query"
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		logFields := logrus.Fields{
			"correlationId": response.DefaultResponseDTO.CorrelationId,
			"userID":        userID,
		}

		timeOut := time.Second * 10
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		hierarchyMongoService := hierarchy.NewHierarchyOperator(mongodb.Operator)
		hierarchyService := hierarchy2.NewHierarchyManager(hierarchy.NewHierarchyOperator(mongodb.Operator))

		mappings, err := hierarchyMongoService.GetAllUserMappingsByQuery(ctx, bson.M{"userId": userID})
		if err != nil {
			logrus.WithFields(logFields).WithError(err).Error("failed to get user mappings")
			response.DefaultResponseDTO.Message = "Failed to retrieve user mappings: " + err.Error()
			response.Status = string(utils.ERROR)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		if len(mappings) == 0 {
			response.DefaultResponseDTO.Message = "No mappings found for the user"
			response.Status = string(utils.SUCCESS)
			c.JSON(http.StatusOK, response)
			return
		}

		var defaultMapping user.Identifiers
		var allMappings []user.Identifiers

		for _, umap := range mappings {
			account, err := hierarchyService.ReadAccount(ctx, umap.AccountID)
			if err != nil {
				logrus.WithFields(logFields).WithError(err).Error("failed to retrieve account details")
				response.DefaultResponseDTO.Message = "Failed to retrieve account details: " + err.Error()
				response.Status = string(utils.ERROR)
				c.JSON(http.StatusInternalServerError, response)
				return
			}

			org, err := hierarchyService.ReadOrganization(ctx, umap.OrgID)
			if err != nil {
				logrus.WithFields(logFields).WithError(err).Error("failed to retrieve organization details")
				response.DefaultResponseDTO.Message = "Failed to retrieve organization details: " + err.Error()
				response.Status = string(utils.ERROR)
				c.JSON(http.StatusInternalServerError, response)
				return
			}

			project, err := hierarchyService.ReadProject(ctx, umap.ProjectID)
			if err != nil {
				logrus.WithFields(logFields).WithError(err).Error("failed to retrieve project details")
				response.DefaultResponseDTO.Message = "Failed to retrieve project details: " + err.Error()
				response.Status = string(utils.ERROR)
				c.JSON(http.StatusInternalServerError, response)
				return
			}

			mappingIdentifiers := user.Identifiers{
				AccountID:   umap.AccountID,
				AccountName: account.Name,
				OrgID:       umap.OrgID,
				OrgName:     org.Name,
				ProjectID:   umap.ProjectID,
				ProjectName: project.Name,
			}

			allMappings = append(allMappings, mappingIdentifiers)

			if umap.IsDefault {
				defaultMapping = mappingIdentifiers
			}
		}

		response.Data = user.UserMapping{
			DefaultMapping: defaultMapping,
			Mappings:       allMappings,
		}

		response.DefaultResponseDTO.Message = "User mappings retrieved successfully"
		response.Status = string(utils.SUCCESS)
		c.JSON(http.StatusOK, response)
	}
}
