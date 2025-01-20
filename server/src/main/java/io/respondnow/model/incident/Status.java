package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum Status {
  STARTED("Started"),
  ACKNOWLEDGED("Acknowledged"),
  INVESTIGATING("Investigating"),
  IDENTIFIED("Identified"),
  MITIGATED("Mitigated"),
  RESOLVED("Resolved");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  Status(String value) {
    this.value = value;
  }
}
