package io.respondnow.repository;

import io.respondnow.model.hierarchy.Organization;
import java.util.Optional;
import javax.validation.constraints.NotBlank;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface OrganizationRepository extends MongoRepository<Organization, String> {
  Optional<Organization> findByOrgIdentifier(@NotBlank String orgIdentifier);
}
