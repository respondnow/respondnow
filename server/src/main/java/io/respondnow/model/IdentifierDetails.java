package io.respondnow.model;

import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class IdentifierDetails {
  private String accountIdentifier;
  private String orgIdentifier;
  private String projectIdentifier;
}
