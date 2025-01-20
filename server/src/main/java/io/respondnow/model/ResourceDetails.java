package io.respondnow.model;

import java.util.List;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class ResourceDetails {
  private String name;
  private String identifier;
  private String description;
  private List<String> tags;
}
