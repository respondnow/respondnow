package io.respondnow.dto.auth;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotEmpty;
import lombok.Getter;
import lombok.Setter;

@Setter
@Getter
public class ChangePasswordInput {

  @NotEmpty(message = "Email cannot be empty")
  @Email(message = "Email should be valid")
  private String email;

  @NotEmpty(message = "Password cannot be empty")
  private String password;

  @NotEmpty(message = "New Password cannot be empty")
  private String newPassword;

  // Getters and Setters
}
