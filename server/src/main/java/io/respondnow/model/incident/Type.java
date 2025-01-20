package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum Type {
  Availability("Availability"),
  Latency("Latency"),
  Security("Security"),
  Other("Other");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  Type(String value) {
    this.value = value;
  }
}
