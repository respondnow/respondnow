package io.respondnow.dto.auth;

import java.util.List;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
public class UserMappingData {
  private UserIdentifiers defaultMapping;
  private List<UserIdentifiers> mappings;
}
