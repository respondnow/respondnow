package io.respondnow.service.hierarchy;

import io.respondnow.dto.auth.UserMappingData;
import io.respondnow.model.hierarchy.UserMapping;
import java.util.List;
import java.util.Optional;

public interface UserMappingService {
  List<UserMapping> findAll();

  Optional<UserMapping> findById(String id);

  UserMapping save(UserMapping userMapping);

  void deleteById(String id);

  UserMapping createUserMapping(
      String userId,
      String accountIdentifier,
      String orgIdentifier,
      String projectIdentifier,
      boolean isDefault);

  UserMappingData getUserMappings(String correlationId, String userId);
}
