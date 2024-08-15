package incident

import (
	"github.com/respondnow/respond/server/pkg/api"
	"github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	"github.com/respondnow/respond/server/utils"
)

type ListResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     ListResponse `json:"data"`
}

type ListResponse struct {
	Content    []incident.Incident `json:"content"`
	Pagination api.Pagination      `json:"pagination"`
}
