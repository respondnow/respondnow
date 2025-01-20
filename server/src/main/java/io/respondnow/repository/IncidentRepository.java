package io.respondnow.repository;

import io.respondnow.model.incident.Incident;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface IncidentRepository extends MongoRepository<Incident, String> {
  // Custom query methods can be added here
}
