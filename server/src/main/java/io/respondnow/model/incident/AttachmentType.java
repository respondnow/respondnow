package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum AttachmentType {
  LINK("Link");

  // Getter method to retrieve the value
  private final String value; // Declare the 'value' field

  // Constructor to set the 'value' field
  AttachmentType(String value) {
    this.value = value;
  }
}
