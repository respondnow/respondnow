package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum RoleType {
  Incident_Commander("Incident_Commander"),
  Communications_Lead("Communications_Lead");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  RoleType(String value) {
    this.value = value;
  }
  @Override public String toString() {
    return value;
  }
}

