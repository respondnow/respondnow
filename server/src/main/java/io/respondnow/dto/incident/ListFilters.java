package io.respondnow.dto.incident;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.respondnow.model.incident.IncidentChannelType;
import io.respondnow.model.incident.Severity;
import io.respondnow.model.incident.Status;
import io.respondnow.model.incident.Type;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class ListFilters {

  private Type type;
  private Severity severity;

  @JsonProperty("incidentChannelType")
  private IncidentChannelType incidentChannelType;

  private Status status;
  private String active;
}
