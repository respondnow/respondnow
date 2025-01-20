package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum Severity {
  SEV0("SEV0"),
  SEV1("SEV1"),
  SEV2("SEV2");

  private final String description;

  // Constructor to set the description field
  Severity(String description) {
    this.description = description;
  }

  // Static method to map a string to the corresponding enum
  public static Severity fromString(String severityStr) {
    for (Severity severity : Severity.values()) {
      if (severity.name().equals(severityStr)) {
        return severity;
      }
    }
    return SEV2; // default to SEV2 if not found
  }
}
