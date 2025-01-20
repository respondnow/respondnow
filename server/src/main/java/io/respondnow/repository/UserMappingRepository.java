package io.respondnow.repository;

import io.respondnow.model.hierarchy.UserMapping;
import jakarta.validation.constraints.NotBlank;
import java.util.List;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface UserMappingRepository extends MongoRepository<UserMapping, String> {
  List<UserMapping> findByUserId(@NotBlank String userId);
}
