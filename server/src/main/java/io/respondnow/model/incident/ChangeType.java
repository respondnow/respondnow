package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum ChangeType {
  SEVERITY("Severity"),
  STATUS("Status"),
  COMMENT("Comment"),
  SUMMARY("Summary"),
  ROLES("Roles"),
  SLACK_CHANNEL_CREATED("Slack Channel Created"),
  INCIDENT_CREATED("Incident Created");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  ChangeType(String value) {
    this.value = value;
  }
}
