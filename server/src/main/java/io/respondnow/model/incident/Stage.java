package io.respondnow.model.incident;

import io.respondnow.model.user.UserDetails;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class Stage {

  private String stageId;
  private Status type;
  private Long duration;
  private Long createdAt;
  private Long updatedAt;
  private UserDetails userDetails;
}
