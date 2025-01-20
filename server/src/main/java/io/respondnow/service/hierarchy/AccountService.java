package io.respondnow.service.hierarchy;

import io.respondnow.model.hierarchy.Account;
import java.util.List;
import java.util.Optional;

public interface AccountService {
  List<Account> findAll();

  Optional<Account> findById(String id);

  Account save(Account account);

  void deleteById(String id);

  Account createAccount(Account account);

  Account createAccountWithRetry(Account account);

  void deleteAccount(String accountIdentifier);

  Account getAccount(String accountIdentifier);

  Iterable<Account> getAllAccounts();
}
