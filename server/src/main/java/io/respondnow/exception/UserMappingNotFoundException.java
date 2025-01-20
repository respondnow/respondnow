package io.respondnow.exception;

public class UserMappingNotFoundException extends RuntimeException {
  public UserMappingNotFoundException(String message) {
    super(message);
  }
}
