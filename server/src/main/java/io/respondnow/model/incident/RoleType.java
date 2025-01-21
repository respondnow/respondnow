package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum RoleType {
  Incident_Commander("Incident_Commander"),
  Communications_Lead("Communications_Lead");


  private final String value;
  private final String displayValue;

  // Constructor to set the 'value' field
  RoleType(String value) {
    this.value = value;
    this.displayValue = value.replace("_", " ");
  }

  @Override public String toString() {
    return value;
  }
}

