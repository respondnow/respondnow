package io.respondnow.dto.auth;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class UserMappingResponseDTO {
  private String correlationId;
  private String message;
  private String status;
  private UserMappingData data;

  // Constructor
  public UserMappingResponseDTO(
      String correlationId, String message, String status, UserMappingData data) {
    this.status = status;
    this.message = message;
    this.correlationId = correlationId;
    this.data = data;
  }
}
