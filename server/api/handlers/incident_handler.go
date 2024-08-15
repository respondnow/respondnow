package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/respondnow/respond/server/pkg/constant"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	incidentdb "github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	"github.com/respondnow/respond/server/pkg/incident"
	"github.com/respondnow/respond/server/utils"
	"github.com/sirupsen/logrus"
)

// ListIncidents godoc
//
//	@Summary		List incidents
//	@Description	List incidents
//	@id				ListIncidents
//
//	@Tags			Incident Management
//	@Accept			json
//	@Produce		json
//	@Param			searchTerm		query		string	false	"searchTerm"	"search based on name and tags"
//	@Param			page			query		int		false	"page"			"Default: 0"	@Default(0)
//	@Param			limit			query		int		false	"limit"			"Default: 10"	@Default(10)
//	@Param			correlationId	query		string	false	"correlationId"	"correlationId is used to debug micro svc communication"
//	@Param			all				query		boolean	false	"all"			"get all"	@Default(false)
//	@Success		200				{object}	incident.IncidentList
//	@Failure		400				{object}	utils.DefaultResponseDTO
//	@Failure		404				{object}	utils.DefaultResponseDTO
//	@Failure		500				{object}	utils.DefaultResponseDTO
//	@Router			/incident/list [get]
//
// ListIncidents handles the retrieval of all incidents.
func ListIncidents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response = incident.ListResponseDTO{
			DefaultResponseDTO: utils.DefaultResponseDTO{
				CorrelationId: c.Query(constant.CorrelationID),
			},
		}

		search := c.Query(constant.Search)
		page, err := strconv.ParseInt(c.Query(constant.Page), 10, 64)
		if err != nil {
			page = constant.DefaultPage
		}
		limit, err := strconv.ParseInt(c.Query(constant.Limit), 10, 64)
		if err != nil {
			limit = constant.DefaultLimit
		}
		if page <= 0 {
			page = constant.DefaultPage
		}
		if limit <= 0 {
			limit = constant.DefaultLimit
		}
		if limit > 50 {
			limit = 50
		}

		response.Status = string(utils.SUCCESS)
		result, err := incident.NewIncidentService(
			incidentdb.NewIncidentOperator(mongodb.Operator)).
			ListIncidents(context.TODO(), page, limit, search)
		if err != nil {
			logrus.Errorf("unable to list incidents, error : %v", err)
			response.Status = string(utils.ERROR)
			response.Message = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
		}

		response.Data = incident.ListResponse{
			Content:    result.Data,
			Pagination: result.Pagination,
		}

		c.JSON(http.StatusOK, response)
	}
}
