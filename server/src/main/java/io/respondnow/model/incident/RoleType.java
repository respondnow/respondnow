package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum RoleType {
  INCIDENT_COMMANDER("Incident Commander"),
  COMMUNICATIONS_LEAD("Communications Lead");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  RoleType(String value) {
    this.value = value;
  }
}
