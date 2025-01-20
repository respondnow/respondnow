package io.respondnow.repository;

import io.respondnow.model.user.User;
import java.util.Optional;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface UserRepository extends MongoRepository<User, String> {
  // Method to find a user by their email
  Optional<User> findByEmail(String email);

  // Check if email exists in the database
  boolean existsByEmail(String email);
}
