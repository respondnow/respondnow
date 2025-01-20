package io.respondnow.model.incident;

import lombok.Getter;
import lombok.NoArgsConstructor;

@Getter
@NoArgsConstructor
public enum Severity {
  SEV0("SEV0 - Critical, High Impact"),
  SEV1("SEV1 - Major, Significant Impact"),
  SEV2("SEV2 - Minor, Low Impact");

  private String description;

  Severity(String description) {
    this.description = description;
  }

  public static Severity fromString(String severityStr) {
    for (Severity severity : Severity.values()) {
      if (severity.name().equals(severityStr)) {
        return severity;
      }
    }
    return SEV2; // default to SEV2 if not found
  }
}
