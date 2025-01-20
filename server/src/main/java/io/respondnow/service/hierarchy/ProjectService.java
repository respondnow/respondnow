package io.respondnow.service.hierarchy;

import io.respondnow.model.hierarchy.Project;
import java.util.List;
import java.util.Optional;

public interface ProjectService {
  List<Project> findAll();

  Optional<Project> findById(String id);

  Project save(Project project);

  void deleteById(String id);

  Project createProject(Project project);

  Project createProjectWithRetry(Project project);

  void deleteProject(String projectIdentifier);

  Project getProject(String projectIdentifier);

  Iterable<Project> getAllProjects();
}
