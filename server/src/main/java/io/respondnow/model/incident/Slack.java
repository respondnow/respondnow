package io.respondnow.model.incident;

import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
public class Slack {

  private String teamId;
  private String teamName;
  private String teamDomain;
  private String channelId;
  private String channelName;
  private String channelReference;
  private String channelDescription;
  private ChannelStatus channelStatus;
}
