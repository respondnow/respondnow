package io.respondnow.model.incident;

import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class Attachment {

  private AttachmentType type;
  private String description;
  private String url;
}
