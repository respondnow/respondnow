package io.respondnow.service.hierarchy;

import io.respondnow.model.hierarchy.Organization;
import java.util.List;
import java.util.Optional;

public interface OrganizationService {
  List<Organization> findAll();

  Optional<Organization> findById(String id);

  Organization save(Organization organization);

  void deleteById(String id);

  Organization createOrganization(Organization organization);

  Organization createOrganizationWithRetry(Organization organization);

  void deleteOrganization(String organizationId);

  Organization getOrganization(String organizationId);

  Iterable<Organization> getAllOrganizations();
}
