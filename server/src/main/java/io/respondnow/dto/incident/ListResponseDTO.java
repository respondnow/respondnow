package io.respondnow.dto.incident;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.respondnow.dto.DefaultResponseDTO;
import io.respondnow.model.api.Pagination;
import io.respondnow.model.incident.Incident;
import java.util.List;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@JsonInclude(JsonInclude.Include.NON_NULL)
@Builder
public class ListResponseDTO extends DefaultResponseDTO {

  @JsonProperty("data")
  private ListResponse listResponse;

  @Getter
  @Setter
  @JsonInclude(JsonInclude.Include.NON_NULL)
  @Builder
  public static class ListResponse {
    private List<Incident> content;
    private Pagination pagination;
    private String correlationID;
  }
}
