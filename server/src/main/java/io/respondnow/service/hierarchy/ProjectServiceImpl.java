package io.respondnow.service.hierarchy;

import io.respondnow.exception.projectIdentifierAlreadyExistsException;
import io.respondnow.exception.ProjectNotFoundException;
import io.respondnow.model.hierarchy.Project;
import io.respondnow.repository.ProjectRepository;
import java.util.List;
import java.util.Optional;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class ProjectServiceImpl implements ProjectService {
  @Autowired private ProjectRepository projectRepository;

  public List<Project> findAll() {
    return projectRepository.findAll();
  }

  public Optional<Project> findById(String id) {
    return projectRepository.findById(id);
  }

  public Project save(Project project) {
    return projectRepository.save(project);
  }

  public void deleteById(String id) {
    projectRepository.deleteById(id);
  }

  @Transactional
  public Project createProject(Project project) {
    // Check if the project already exists
    Optional<Project> existingProject = projectRepository.findByProjectIdentifier(project.getProjectIdentifier());
    if (existingProject.isPresent()) {
      throw new projectIdentifierAlreadyExistsException("Project with the given id already exists");
    }
    return projectRepository.save(project);
  }

  @Retryable(
      value = Exception.class,
      maxAttempts = 3,
      backoff = @Backoff(delay = 2000, multiplier = 1.5))
  public Project createProjectWithRetry(Project project) {
    // Check if the project already exists
    Optional<Project> existingProject = projectRepository.findByProjectIdentifier(project.getProjectIdentifier());
    if (existingProject.isPresent()) {
      throw new projectIdentifierAlreadyExistsException("Project with the given id already exists");
    }
    return projectRepository.save(project);
  }

  @Transactional
  public void deleteProject(String projectIdentifier) {
    Project project =
        projectRepository
            .findByProjectIdentifier(projectIdentifier)
            .orElseThrow(() -> new ProjectNotFoundException("Project not found"));
    project.setRemoved(true); // Soft delete the project
    projectRepository.save(project);
  }

  @Transactional
  public Project getProject(String projectIdentifier) {
    return projectRepository
        .findByProjectIdentifier(projectIdentifier)
        .orElseThrow(() -> new ProjectNotFoundException("Project not found"));
  }

  public Iterable<Project> getAllProjects() {
    return projectRepository.findAll();
  }
}
