package io.respondnow.model.incident;

import javax.validation.constraints.NotNull;

import io.respondnow.model.user.UserDetails;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class Timeline {

  private String id;

  @NotNull private ChangeType type;

  private Long createdAt;

  private Long updatedAt;

  private String previousState;

  private String currentState;

  private Slack slack;

  private UserDetails userDetails;

  private String message;

  private Object additionalDetails;
}
