package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum ChannelSource {
  SLACK("Slack");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  ChannelSource(String value) {
    this.value = value;
  }
}
