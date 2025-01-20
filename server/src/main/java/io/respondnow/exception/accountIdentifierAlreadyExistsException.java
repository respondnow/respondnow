package io.respondnow.exception;

public class accountIdentifierAlreadyExistsException extends RuntimeException {
  public accountIdentifierAlreadyExistsException(String message) {
    super(message);
  }
}
