package io.respondnow.model.incident;

import io.respondnow.model.user.UserDetails;
import java.util.ArrayList;
import java.util.List;
import javax.validation.constraints.NotNull;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

@Data
@NoArgsConstructor
@Getter
@Setter
@Document(collection = "incidents")
public class Incident {

  @Id private String id;

  private String accountIdentifier;
  private String orgIdentifier;
  private String projectIdentifier;
  private String name;
  private String identifier;
  private String description;
  private List<String> tags;
  @NotNull private Type type;

  @NotNull private Severity severity;

  @NotNull private Status status;

  @NotNull private String summary;

  @NotNull private String comment;

  @NotNull private Boolean active;

  private List<Service> services;
  private List<Environment> environments;
  private List<Functionality> functionalities;
  private List<Role> roles;
  private List<Timeline> timelines;
  private List<Stage> stages;
  private List<Channel> channels;
  private IncidentChannel incidentChannel;
  private List<Conference> conferenceDetails;
  private List<Attachment> attachments;
  private Long createdAt;
  private Long updatedAt;
  private UserDetails createdBy;
  private UserDetails updatedBy;
  private Long removedAt;
  private Boolean removed;

  /**
   * Adds a Timeline entry to the Incident's timeline history.
   *
   * @param timeline The Timeline object to be added.
   */
  public void addTimeline(Timeline timeline) {
    if (this.timelines == null) {
      this.timelines = new ArrayList<>();
    }
    this.timelines.add(timeline);
  }

  /**
   * Optionally, you can add a method to remove a Timeline entry.
   *
   * @param timeline The Timeline object to be removed.
   */
  public void removeTimeline(Timeline timeline) {
    if (this.timelines != null) {
      this.timelines.remove(timeline);
    }
  }
}
