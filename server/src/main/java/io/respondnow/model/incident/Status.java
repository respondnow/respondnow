package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum Status {
  Started("Started"),
  Acknowledged("Acknowledged"),
  Investigating("Investigating"),
  Identified("Identified"),
  Mitigated("Mitigated"),
  Resolved("Resolved");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  Status(String value) {
    this.value = value;
  }
}
