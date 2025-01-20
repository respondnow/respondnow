package io.respondnow.service.auth;

import io.respondnow.dto.auth.AddUserInput;
import io.respondnow.dto.auth.ChangePasswordInput;
import io.respondnow.dto.auth.LoginUserInput;
import io.respondnow.exception.EmailAlreadyExistsException;
import io.respondnow.exception.UserNotFoundException;
import io.respondnow.model.user.User;
import io.respondnow.repository.UserRepository;
import io.respondnow.util.JWTUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class AuthServiceImpl implements AuthService {

  @Autowired private UserRepository userRepository;

  @Autowired private BCryptPasswordEncoder passwordEncoder;

  @Autowired private JWTUtil jwtUtil;

  @Override
  public User login(LoginUserInput input) {
    User user =
        userRepository
            .findByEmail(input.getEmail())
            .orElseThrow(() -> new UserNotFoundException("User not found"));

    if (!passwordEncoder.matches(input.getPassword(), user.getPassword())) {
      throw new UserNotFoundException("Invalid credentials");
    }
    user.setLastLoginAt(System.currentTimeMillis());
    return userRepository.save(user);
  }

  @Override
  public User changePassword(ChangePasswordInput input) {
    User user =
        userRepository
            .findByEmail(input.getEmail())
            .orElseThrow(() -> new UserNotFoundException("User not found"));

    user.setPassword(passwordEncoder.encode(input.getNewPassword()));
    user.setChangePasswordRequired(false);
    user.setActive(true);
    user.setUpdatedAt(System.currentTimeMillis());
    return userRepository.save(user);
  }

  @Override
  public User signup(AddUserInput input) {
    // Check if the email already exists
    if (userRepository.existsByEmail(input.getEmail())) {
      throw new EmailAlreadyExistsException("User email already exists");
    }
    User user = new User();
    user.setEmail(input.getEmail());
    user.setPassword(passwordEncoder.encode(input.getPassword()));
    user.setName(input.getName());
    user.setUserId(input.getUserId());
    user.setActive(false);
    user.setChangePasswordRequired(true);
    user.setCreatedAt(System.currentTimeMillis());
    user.setUpdatedAt(System.currentTimeMillis());
    user.setRemoved(false);
    return userRepository.save(user);
  }

  @Retryable(
      value = Exception.class,
      maxAttempts = 3,
      backoff = @Backoff(delay = 2000, multiplier = 1.5))
  public User signupWithRetry(AddUserInput input) {
    // Check if the email already exists
    if (userRepository.existsByEmail(input.getEmail())) {
      throw new EmailAlreadyExistsException("User email already exists");
    }
    User user = new User();
    user.setEmail(input.getEmail());
    user.setPassword(passwordEncoder.encode(input.getPassword()));
    user.setName(input.getName());
    user.setUserId(input.getUserId());
    user.setActive(false);
    user.setChangePasswordRequired(true);
    user.setCreatedAt(System.currentTimeMillis());
    user.setUpdatedAt(System.currentTimeMillis());
    user.setRemoved(false);
    return userRepository.save(user);
  }
}
