package io.respondnow.model.incident;

import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
public class Channel {

  private String id;
  private String teamId;
  private String name;
  private ChannelSource source;
  private String url;
  private ChannelStatus status;
}
