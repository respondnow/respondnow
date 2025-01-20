package io.respondnow.dto.incident;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.respondnow.model.IdentifierDetails;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class GetFilters extends IdentifierDetails {

  @JsonProperty("incidentId")
  private String incidentId;
}
