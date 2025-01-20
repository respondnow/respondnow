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
public class CreateResponseDTO extends DefaultResponseDTO {

  @JsonProperty("data")
  private CreateResponse createResponse;

  @Getter
  @Setter
  @JsonInclude(JsonInclude.Include.NON_NULL)
  public static class CreateResponse {
    private Incident incident;
    private String correlationID;
  }
}
