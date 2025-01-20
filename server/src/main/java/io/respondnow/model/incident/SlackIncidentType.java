package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum SlackIncidentType {
  CLOSED("Closed"),
  OPEN("Open");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  SlackIncidentType(String value) {
    this.value = value;
  }
}
