package io.respondnow.dto.auth;

import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@Getter
@Setter
@NoArgsConstructor
public class UserIdentifiers {
  private String accountIdentifier;
  private String accountName;
  private String orgIdentifier;
  private String orgName;
  private String projectIdentifier;
  private String projectName;
}
