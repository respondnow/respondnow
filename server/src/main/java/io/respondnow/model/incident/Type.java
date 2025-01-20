package io.respondnow.model.incident;

import lombok.Getter;

@Getter
public enum Type {
  Availability("Availability"),
  Latency("Latency"),
  Security("Security"),
  Other("Other");

  private final String name;

  Type(String value) {
    this.name = value;
  }
}
