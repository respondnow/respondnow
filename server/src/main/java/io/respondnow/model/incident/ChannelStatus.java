package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum ChannelStatus {
  Operational("Operational");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  ChannelStatus(String value) {
    this.value = value;
  }
}
