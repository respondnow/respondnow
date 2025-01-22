package io.respondnow.service.hierarchy;

import io.respondnow.exception.AccountNotFoundException;
import io.respondnow.exception.accountIdentifierAlreadyExistsException;
import io.respondnow.model.hierarchy.Account;
import io.respondnow.repository.AccountRepository;
import java.util.List;
import java.util.Optional;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class AccountServiceImpl implements AccountService {
  @Autowired private AccountRepository accountRepository;

  public List<Account> findAll() {
    return accountRepository.findAll();
  }

  public Optional<Account> findById(String id) {
    return accountRepository.findById(id);
  }

  public Account save(Account account) {
    return accountRepository.save(account);
  }

  public void deleteById(String id) {
    accountRepository.deleteById(id);
  }

  @Transactional
  public Account createAccount(Account account) {
    // Check if the account already exists
    Optional<Account> existingAccount =
        accountRepository.findByAccountIdentifier(account.getAccountIdentifier());
    if (existingAccount.isPresent()) {
      throw new accountIdentifierAlreadyExistsException(
          "Account with the given account_id already exists");
    }
    return accountRepository.save(account);
  }

  @Retryable(
      value = Exception.class,
      maxAttempts = 3,
      backoff = @Backoff(delay = 2000, multiplier = 1.5))
  public Account createAccountWithRetry(Account account) {
    // Check if the account already exists
    Optional<Account> existingAccount =
        accountRepository.findByAccountIdentifier(account.getAccountIdentifier());
    if (existingAccount.isPresent()) {
      throw new accountIdentifierAlreadyExistsException(
          "Account with the given account_id already exists");
    }
    return accountRepository.save(account);
  }

  @Transactional
  public void deleteAccount(String accountIdentifier) {
    Account account =
        accountRepository
            .findByAccountIdentifier(accountIdentifier)
            .orElseThrow(() -> new AccountNotFoundException("Account not found"));

    account.setRemoved(true); // Soft delete the account
    accountRepository.save(account);
  }

  @Transactional
  public Account getAccount(String accountIdentifier) {
    return accountRepository
        .findByAccountIdentifier(accountIdentifier)
        .orElseThrow(() -> new AccountNotFoundException("Account not found"));
  }

  public Iterable<Account> getAllAccounts() {
    return accountRepository.findAll();
  }
}
