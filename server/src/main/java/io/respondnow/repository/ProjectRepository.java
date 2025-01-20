package io.respondnow.repository;

import io.respondnow.model.hierarchy.Project;
import java.util.Optional;
import javax.validation.constraints.NotBlank;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface ProjectRepository extends MongoRepository<Project, String> {
  Optional<Project> findByProjectIdentifier(@NotBlank String projectIdentifier);
}
