package io.respondnow.service.hierarchy;

import io.respondnow.exception.OrganizationNotFoundException;
import io.respondnow.exception.orgIdentifierAlreadyExistsException;
import io.respondnow.model.hierarchy.Organization;
import io.respondnow.repository.OrganizationRepository;
import java.util.List;
import java.util.Optional;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class OrganizationServiceImpl implements OrganizationService {
  @Autowired private OrganizationRepository organizationRepository;

  public List<Organization> findAll() {
    return organizationRepository.findAll();
  }

  public Optional<Organization> findById(String id) {
    return organizationRepository.findById(id);
  }

  public Organization save(Organization organization) {
    return organizationRepository.save(organization);
  }

  public void deleteById(String id) {
    organizationRepository.deleteById(id);
  }

  @Transactional
  public Organization createOrganization(Organization organization) {
    // Check if the org already exists
    Optional<Organization> existingOrg =
        organizationRepository.findByOrgIdentifier(organization.getOrgIdentifier());
    if (existingOrg.isPresent()) {
      throw new orgIdentifierAlreadyExistsException(
          "Organization with the given org_id already exists");
    }
    return organizationRepository.save(organization);
  }

  @Retryable(
      value = Exception.class,
      maxAttempts = 3,
      backoff = @Backoff(delay = 2000, multiplier = 1.5))
  public Organization createOrganizationWithRetry(Organization organization) {
    // Check if the org already exists
    Optional<Organization> existingOrg =
        organizationRepository.findByOrgIdentifier(organization.getOrgIdentifier());
    if (existingOrg.isPresent()) {
      throw new orgIdentifierAlreadyExistsException(
          "Organization with the given org_id already exists");
    }
    return organizationRepository.save(organization);
  }

  @Transactional
  public void deleteOrganization(String organizationId) {
    Organization org =
        organizationRepository
            .findByOrgIdentifier(organizationId)
            .orElseThrow(() -> new OrganizationNotFoundException("Organization not found"));
    org.setRemoved(true); // Soft delete the organization
    organizationRepository.save(org);
  }

  @Transactional
  public Organization getOrganization(String organizationId) {
    return organizationRepository
        .findByOrgIdentifier(organizationId)
        .orElseThrow(() -> new OrganizationNotFoundException("Organization not found"));
  }

  public Iterable<Organization> getAllOrganizations() {
    return organizationRepository.findAll();
  }
}
