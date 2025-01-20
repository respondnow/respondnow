package io.respondnow.model.incident;

import io.respondnow.model.user.UserDetails;
import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
public class Role {

  private RoleType roleType;
  private UserDetails userDetails;
}
