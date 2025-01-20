package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum ChangeType {
  Severity("Severity"),
  Status("Status"),
  Comment("Comment"),
  Summary("Summary"),
  Roles("Roles"),
  Slack_Channel_Created("Slack_Channel_Created"),
  Incident_Created("Incident_Created");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  ChangeType(String value) {
    this.value = value;
  }
}
