package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/respondnow/respond/server/api/middleware"
	"github.com/respondnow/respond/server/pkg/constant"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	incidentdb "github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	"github.com/respondnow/respond/server/pkg/incident"
	"github.com/respondnow/respond/server/utils"
	"github.com/sirupsen/logrus"
)

// CreateIncident godoc
//
//	@Summary		Create an incident
//	@Description	Create an incident
//	@id				CreateIncident
//	@Security		ApiKeyAuth
//
//	@Tags			Incident
//	@Accept			json
//	@Produce		json
//	@Param			"create incident"	body		incident.CreateRequest	true	"Create an incident"
//	@Param			accountIdentifier	query		string					true	"accountIdentifier"	"accountIdentifier is the account where you want to access the resource"
//	@Param			orgIdentifier		query		string					false	"orgIdentifier"		"orgIdentifier is the org where you want to access the resource"
//	@Param			projectIdentifier	query		string					false	"projectIdentifier"	"projectIdentifier is the project where you want to access the resource"
//	@Param			correlationId		query		string					false	"correlationId"
//	@Success		200					{object}	incident.CreateResponseDTO
//	@Failure		400					{object}	utils.DefaultResponseDTO
//	@Failure		404					{object}	utils.DefaultResponseDTO
//	@Failure		500					{object}	utils.DefaultResponseDTO
//	@Router			/incident/create [post]
func CreateIncident() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = incident.CreateResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query("correlationId"),
			},
		}

		if response.DefaultResponseDTO.CorrelationId == "" {
			randLen := 16
			response.DefaultResponseDTO.CorrelationId = utils.NewUtils().RandStringBytes(randLen)
		}

		accountId := c.Query("accountIdentifier")
		orgId := c.Query("orgIdentifier")
		projectId := c.Query("projectIdentifier")

		if accountId == "" {
			response.Status = string(utils.ERROR)
			response.Message = "accountIdentifier is required"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		var payload incident.CreateRequest
		if err := c.BindJSON(&payload); err != nil {
			logrus.WithField("correlationId", response.DefaultResponseDTO.CorrelationId).WithError(err).Error("failed to bind request payload")
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		if err := validator.New().Struct(payload); err != nil {
			logrus.WithField("correlationId", response.DefaultResponseDTO.CorrelationId).WithError(err).Error("failed to validate the request")
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		incident, err := incident.NewIncidentService(incidentdb.NewIncidentOperator(mongodb.Operator),
			accountId, orgId, projectId).
			Create(context.TODO(), payload, middleware.CurrentUser{}, response.CorrelationId)
		if err != nil {
			logrus.WithField("correlationId", response.DefaultResponseDTO.CorrelationId).WithError(err).Error("failed to create incident")
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		response.Status = string(utils.SUCCESS)
		response.Data = incident
		c.JSON(http.StatusOK, response)
	}
}

// ListIncidents godoc
//
//	@Summary		List incidents
//	@Description	List incidents
//	@id				ListIncidents
//
//	@Tags			Incident Management
//	@Accept			json
//	@Produce		json
//	@Param			accountIdentifier	query		string							true	"accountIdentifier"		"accountIdentifier is the account where you want to access the resource"
//	@Param			orgIdentifier		query		string							false	"orgIdentifier"			"orgIdentifier is the org where you want to access the resource"
//	@Param			projectIdentifier	query		string							false	"projectIdentifier"		"projectIdentifier is the project where you want to access the resource"
//	@Param			type				query		incidentdb.Type					false	"type"					"type of the incident"
//	@Param			severity			query		incidentdb.Severity				false	"severity"				"severity of the incident"
//	@Param			status				query		incidentdb.Status				false	"status"				"status of the incident"
//	@Param			active				query		bool							false	"active"				"whether incident is active or not"
//	@Param			incidentChannelType	query		incidentdb.IncidentChannelType	false	"incidentChannelType"	"type of the incident channel"
//	@Param			search				query		string							false	"search"				"search based on name and tags"
//	@Param			page				query		int								false	"page"					"Default: 0"	@Default(0)
//	@Param			limit				query		int								false	"limit"					"Default: 10"	@Default(10)
//	@Param			correlationId		query		string							false	"correlationId"			"correlationId is used to debug micro svc communication"
//	@Param			all					query		boolean							false	"all"					"get all"	@Default(false)
//	@Success		200					{object}	incident.ListResponseDTO
//	@Failure		400					{object}	utils.DefaultResponseDTO
//	@Failure		404					{object}	utils.DefaultResponseDTO
//	@Failure		500					{object}	utils.DefaultResponseDTO
//	@Router			/incident/list [get]
func ListIncidents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = incident.ListResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: utils.NewUtils().GetCorrelationID(c),
			},
		}

		aID := c.Query(constant.AccountIdentifier)
		oID := c.Query(constant.OrgIdentifier)
		pID := c.Query(constant.ProjectIdentifier)
		incidentType := c.Query(constant.Type)
		severity := c.Query(constant.Severity)
		status := c.Query(constant.Status)
		active := c.Query(constant.Active)
		incidentChannelType := c.Query(constant.IncidentChannelType)

		search := c.Query(constant.Search)
		page, limit, all := utils.NewUtils().GetPagination(c)

		response.Status = string(utils.SUCCESS)
		result, err := incident.NewIncidentService(
			incidentdb.NewIncidentOperator(mongodb.Operator), aID, oID, pID).
			List(context.TODO(), c.Request.Header.Get(constant.Authorization),
				incident.ListFilters{
					Type:                incidentdb.Type(incidentType),
					Severity:            incidentdb.Severity(severity),
					IncidentChannelType: incidentdb.IncidentChannelType(incidentChannelType),
					Status:              incidentdb.Status(status),
					Active:              active,
				}, response.CorrelationId, search, limit, page, all)
		if err != nil {
			logrus.Errorf("unable to list incidents, error : %v", err)
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
		}
		response.Data = result

		c.JSON(http.StatusOK, response)
	}
}
