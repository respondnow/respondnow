package io.respondnow.dto.auth;

import io.respondnow.model.user.User;
import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Getter;

@Getter
public class SignupResponseDTO {

  // Getters (No setters needed as fields are final)
  @Schema(description = "Status of the response", example = "success")
  private final String status;

  @Schema(
      description = "Message with additional information about the request",
      example = "User registered successfully")
  private final String message;

  @Schema(description = "JWT token for the newly registered user")
  private final String token;

  @Schema(description = "User data object")
  private final User data;

  // Constructor
  public SignupResponseDTO(String status, String message, String token, User data) {
    this.status = status;
    this.message = message;
    this.token = token;
    this.data = data;
  }
}
