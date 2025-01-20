package io.respondnow.model.user;

import io.respondnow.model.incident.ChannelSource;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class UserDetails {
  private String userId;
  private String userName;
  private String email;
  private String name;
  private ChannelSource source;
}
