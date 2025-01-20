package io.respondnow.model.incident;

import io.respondnow.model.user.UserDetails;
import javax.validation.constraints.NotNull;
import lombok.*;

@Data
@NoArgsConstructor
@Getter
@Setter
@AllArgsConstructor
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
