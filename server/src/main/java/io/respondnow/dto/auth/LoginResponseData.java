package io.respondnow.dto.auth;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class LoginResponseData {
  @Schema(description = "JWT token for the user")
  private String token;

  @Schema(description = "Timestamp of the user's last login", example = "1630421333000")
  private long lastLoginAt;

  @Schema(description = "Indicates if the user needs to change their password", example = "true")
  private boolean changeUserPassword;
}
