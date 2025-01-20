package io.respondnow.dto.incident;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.respondnow.model.ResourceDetails;
import io.respondnow.model.incident.*;
import java.util.List;
import lombok.*;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
public class CreateRequest extends ResourceDetails {

  @JsonProperty("type")
  private Type type;

  @JsonProperty("severity")
  private Severity severity;

  @JsonProperty("summary")
  private String summary;

  @JsonProperty("incidentChannel")
  private IncidentChannel incidentChannel;

  @JsonProperty("status")
  private Status status;

  @JsonProperty("services")
  private List<Service> services;

  @JsonProperty("environments")
  private List<Environment> environments;

  @JsonProperty("functionalities")
  private List<Functionality> functionalities;

  @JsonProperty("channels")
  private List<Channel> channels;

  @JsonProperty("roles")
  private List<Role> roles;

  @JsonProperty("addConference")
  private AddConference addConference;

  @JsonProperty("attachments")
  private List<Attachment> attachments;

  @Getter
  @Setter
  @JsonInclude(JsonInclude.Include.NON_NULL)
  public static class AddConference {
    private String type;
  }
}
