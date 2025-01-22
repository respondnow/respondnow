package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum RoleType {
  Incident_Commander("Incident_Commander"),
  Communications_Lead("Communications_Lead");
  // Add other role types as needed

  private final String value;
  private final String displayValue;

  // Constructor to set the 'value' and 'displayValue' fields
  RoleType(String value) {
    this.value = value;
    this.displayValue = value.replace("_", " ");
  }

  /**
   * Converts a string key to its corresponding RoleType enum.
   *
   * @param key The string representation of the role type.
   * @return The corresponding RoleType enum.
   * @throws IllegalArgumentException If the key does not match any RoleType.
   */
  public static RoleType fromValue(String key) {
    if (key == null) {
      throw new IllegalArgumentException("RoleType key cannot be null.");
    }

    for (RoleType roleType : RoleType.values()) {
      if (roleType.getValue().equalsIgnoreCase(key)) {
        return roleType;
      }
    }

    throw new IllegalArgumentException("Unknown RoleType: " + key);
  }

  @Override
  public String toString() {
    return value;
  }
}
