package io.respondnow.dto.auth;

import lombok.*;

@Data
@Getter
@Setter
@NoArgsConstructor
public class GetUserMappingResponseDTO {
  private UserMappingData data;
  private String correlationId;
  private String message;
  private String status;
}
