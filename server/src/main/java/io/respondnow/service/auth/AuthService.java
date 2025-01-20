package io.respondnow.service.auth;

import io.respondnow.dto.auth.AddUserInput;
import io.respondnow.dto.auth.ChangePasswordInput;
import io.respondnow.dto.auth.LoginUserInput;
import io.respondnow.model.user.User;

public interface AuthService {
  User signup(AddUserInput input);

  User signupWithRetry(AddUserInput input);

  User login(LoginUserInput input);

  User changePassword(ChangePasswordInput input);
}
