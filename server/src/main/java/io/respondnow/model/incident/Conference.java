package io.respondnow.model.incident;

import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class Conference {

  private String conferenceId;
  private String type;
  private String url;
}
