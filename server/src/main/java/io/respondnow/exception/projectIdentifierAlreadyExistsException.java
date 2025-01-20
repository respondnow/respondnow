package io.respondnow.exception;

public class projectIdentifierAlreadyExistsException extends RuntimeException {
  public projectIdentifierAlreadyExistsException(String message) {
    super(message);
  }
}
