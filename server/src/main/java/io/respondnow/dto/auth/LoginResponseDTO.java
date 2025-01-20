package io.respondnow.dto.auth;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Getter;

@Getter
public class LoginResponseDTO {

  // Getters (No setters needed as fields are final)
  @Schema(description = "Status of the login response", example = "success")
  private final String status;

  @Schema(description = "Message with additional information about the login")
  private final String message;

  @Schema(description = "User data object")
  private final LoginResponseData data;

  // Constructor
  public LoginResponseDTO(String status, String message, LoginResponseData data) {
    this.status = status;
    this.message = message;
    this.data = data;
  }
}
