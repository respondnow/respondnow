package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum ConferenceType {
  Zoom("Zoom");

  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  ConferenceType(String value) {
    this.value = value;
  }

  // Getter method to retrieve the value
  public String getValue() {
    return value;
  }
}
