package io.respondnow.dto.incident;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.respondnow.dto.DefaultResponseDTO;
import io.respondnow.model.incident.Incident;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@JsonInclude(JsonInclude.Include.NON_NULL)
public class GetResponseDTO extends DefaultResponseDTO {

  @JsonProperty("data")
  private Incident incident;
}
