package io.respondnow.repository;

import io.respondnow.model.incident.Incident;
import java.util.Optional;
import javax.validation.constraints.NotBlank;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface IncidentRepository extends MongoRepository<Incident, String> {
  Optional<Incident> findByIdentifier(@NotBlank String identifier);
}
