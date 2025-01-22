package io.respondnow.controller;

import io.respondnow.dto.incident.GetResponseDTO;
import io.respondnow.dto.incident.ListResponseDTO;
import io.respondnow.model.api.Pagination;
import io.respondnow.model.incident.Incident;
import io.respondnow.model.incident.Severity;
import io.respondnow.model.incident.Status;
import io.respondnow.model.incident.Type;
import io.respondnow.service.incident.IncidentService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.enums.ParameterIn;
import io.swagger.v3.oas.annotations.tags.Tag;
import java.util.List;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

@Tag(name = "Incident Controller", description = "APIs for Incident Management")
@RestController
@RequestMapping("/incident")
public class IncidentController {

    @Autowired
    private IncidentService incidentService;
    @Operation(summary = "List incidents", description = "List all incidents with optional filters")
    @GetMapping("/list")
    @ResponseStatus(HttpStatus.OK)
    public ListResponseDTO listIncidents(
            @Parameter(
                    name = "accountIdentifier",
                    description = "accountIdentifier is the account where you want to access the resource",
                    in = ParameterIn.QUERY,
                    required = false
            )
            @RequestParam String accountIdentifier,

            @Parameter(
                    name = "orgIdentifier",
                    description = "orgIdentifier is the org where you want to access the resource",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) String orgIdentifier,

            @Parameter(
                    name = "projectIdentifier",
                    description = "projectIdentifier is the project where you want to access the resource",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) String projectIdentifier,

            @Parameter(
                    name = "type",
                    description = "type of the incident",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) Type type,

            @Parameter(
                    name = "severity",
                    description = "severity of the incident",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) Severity severity,

            @Parameter(
                    name = "status",
                    description = "status of the incident",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) Status status,

            @Parameter(
                    name = "active",
                    description = "whether incident is active or not",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) Boolean active,

            @Parameter(
                    name = "incidentChannelType",
                    description = "type of the incident channel",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) String incidentChannelType,

            @Parameter(
                    name = "search",
                    description = "search based on name and/or tags",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) String search,

            @Parameter(
                    name = "page",
                    description = "Pagination page, default: 0",
                    in = ParameterIn.QUERY
            )
            @RequestParam(defaultValue = "0") int page,

            @Parameter(
                    name = "limit",
                    description = "Pagination limit, default: 10",
                    in = ParameterIn.QUERY
            )
            @RequestParam(defaultValue = "10") int limit,

            @Parameter(
                    name = "all",
                    description = "If true, returns all incidents without pagination",
                    in = ParameterIn.QUERY
            )
            @RequestParam(defaultValue = "false") boolean all,

            @Parameter(
                    name = "correlationId",
                    description = "correlationId is used to debug micro service communication",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) String correlationId
    ) {
        Criteria criteria = new Criteria();

        criteria.and("accountIdentifier").is(accountIdentifier);

        if (orgIdentifier != null && !orgIdentifier.isEmpty()) {
            criteria.and("orgIdentifier").is(orgIdentifier);
        }
        if (projectIdentifier != null && !projectIdentifier.isEmpty()) {
            criteria.and("projectIdentifier").is(projectIdentifier);
        }
        if (type != null) {
            criteria.and("type").is(type);
        }
        if (severity != null) {
            criteria.and("severity").is(severity);
        }
        if (status != null) {
            criteria.and("status").is(status);
        }
        if (active != null) {
            criteria.and("active").is(active);
        }
        if (incidentChannelType != null && !incidentChannelType.isEmpty()) {
            criteria.and("incidentChannel.type").is(incidentChannelType);
        }
        if (search != null && !search.isEmpty()) {
            criteria.orOperator(
                    Criteria.where("name").regex(search, "i"),
                    Criteria.where("tags").regex(search, "i")
            );
        }

        Query query = new Query(criteria);

        if (!all) {
            query.skip((long) page * limit);
            query.limit(limit);
        }

        List<Incident> incidents = incidentService.listIncidents(query);
        long total = incidentService.countIncidents(query);

        return ListResponseDTO.builder()
                .listResponse(ListResponseDTO.ListResponse.builder()
                        .content(incidents)
                        .pagination(Pagination.builder().totalItems(total).build())
                        .correlationID(correlationId)
                        .build())
                .build();
    }


    @Operation(summary = "Get incident", description = "Fetch a specific incident by its identifier")
    @GetMapping("/{incidentIdentifier}")
    @ResponseStatus(HttpStatus.OK)
    public GetResponseDTO getIncident(
            @Parameter(
                    name = "incidentIdentifier",
                    description = "Incident identifier (ObjectId)",
                    required = true,
                    in = ParameterIn.PATH
            )
            @PathVariable String incidentIdentifier,

            @Parameter(
                    name = "accountIdentifier",
                    description = "accountIdentifier is the account where you want to access the resource",
                    in = ParameterIn.QUERY,
                    required = true
            )
            @RequestParam String accountIdentifier,

            @Parameter(
                    name = "orgIdentifier",
                    description = "orgIdentifier is the org where you want to access the resource",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) String orgIdentifier,

            @Parameter(
                    name = "projectIdentifier",
                    description = "projectIdentifier is the project where you want to access the resource",
                    in = ParameterIn.QUERY
            )
            @RequestParam(required = false) String projectIdentifier
    ) {
        Incident incident = incidentService.getIncidentById(incidentIdentifier);
        return GetResponseDTO.builder().incident(incident).build();
    }
}
