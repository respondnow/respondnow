package io.respondnow.repository;

import io.respondnow.model.hierarchy.Account;
import java.util.Optional;
import javax.validation.constraints.NotBlank;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface AccountRepository extends MongoRepository<Account, String> {
  Optional<Account> findByAccountIdentifier(@NotBlank String accountIdentifier);
}
