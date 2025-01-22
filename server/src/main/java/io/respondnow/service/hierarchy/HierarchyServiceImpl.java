package io.respondnow.service.hierarchy;

import io.respondnow.dto.auth.AddUserInput;
import io.respondnow.exception.EmailAlreadyExistsException;
import io.respondnow.model.hierarchy.Account;
import io.respondnow.model.hierarchy.Organization;
import io.respondnow.model.hierarchy.Project;
import io.respondnow.model.hierarchy.UserMapping;
import io.respondnow.model.user.User;
import io.respondnow.service.auth.AuthService;
import io.respondnow.service.user.UserService;
import java.util.concurrent.CompletableFuture;
import javax.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Async;
import org.springframework.scheduling.annotation.EnableAsync;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@EnableAsync
public class HierarchyServiceImpl implements HierarchyService {

  private static final Logger logger = LoggerFactory.getLogger(HierarchyServiceImpl.class);
  @Autowired private UserService userService;
  @Autowired private AuthService authService;
  @Autowired private UserMappingService userMappingService;
  @Autowired private AccountService accountService;
  @Autowired private OrganizationService organizationService;
  @Autowired private ProjectService projectService;

  @Value("${hierarchy.defaultAccount.id:default_account_id}")
  private String defaultAccountId;

  @Value("${hierarchy.defaultAccount.name:Default Account}")
  private String defaultAccountName;

  @Value("${hierarchy.defaultOrg.id:default_org_id}")
  private String defaultOrgId;

  @Value("${hierarchy.defaultOrg.name:Default Org}")
  private String defaultOrgName;

  @Value("${hierarchy.defaultProject.id:default_project_id}")
  private String defaultProjectId;

  @Value("${hierarchy.defaultProject.name:Default Project}")
  private String defaultProjectName;

  @Value("${hierarchy.defaultUser.name:Admin}")
  private String defaultUserName;

  @Value("${hierarchy.defaultUser.email:admin@respondnow.io}")
  private String defaultUserEmail;

  @Value("${hierarchy.defaultUser.password:respondnow}")
  private String defaultUserPassword;

  @Async
  @Transactional
  public CompletableFuture<Void> backgroundProcess() {
    try {
      // Create default user
      User createdUser = createDefaultUser();
      if (createdUser == null) {
        return CompletableFuture.completedFuture(null);
      }

      // Create default account
      Account defaultAccount = createDefaultAccount();
      if (defaultAccount == null) {
        fallbackCleanup(createdUser, null, null, null);
        return CompletableFuture.completedFuture(null);
      }

      // Create default organization
      Organization defaultOrg = createDefaultOrganization(defaultAccount);
      if (defaultOrg == null) {
        fallbackCleanup(createdUser, defaultAccount, null, null);
        return CompletableFuture.completedFuture(null);
      }

      // Create default project
      Project defaultProject = createDefaultProject(defaultAccount, defaultOrg);
      if (defaultProject == null) {
        fallbackCleanup(createdUser, defaultAccount, defaultOrg, null);
        return CompletableFuture.completedFuture(null);
      }

      // Create user mapping
      UserMapping defaultUserMapping =
          createUserMapping(createdUser, defaultAccount, defaultOrg, defaultProject);
      if (defaultUserMapping == null) {
        fallbackCleanup(createdUser, defaultAccount, defaultOrg, defaultProject);
        return CompletableFuture.completedFuture(null);
      }

      // Log success
      logger.info("All resources have been successfully created.");
    } catch (EmailAlreadyExistsException e) {
      logger.info("Default user already exists.");
    } catch (Exception e) {
      System.err.println("Error in background process: " + e.getMessage());
    }
    return CompletableFuture.completedFuture(null);
  }

  private User createDefaultUser() {
    AddUserInput input = new AddUserInput();
    input.setUserId("admin");
    input.setName(defaultUserName);
    input.setEmail(defaultUserEmail);
    input.setPassword(defaultUserPassword);
    return authService.signupWithRetry(input);
  }

  private Account createDefaultAccount() {
    Account defaultAccount = new Account();
    defaultAccount.setAccountIdentifier(defaultAccountId);
    defaultAccount.setName(defaultAccountName);
    defaultAccount.setCreatedBy("SYSTEM");
    defaultAccount.setCreatedAt(System.currentTimeMillis());
    defaultAccount.setUpdatedBy("SYSTEM");

    return accountService.createAccountWithRetry(defaultAccount);
  }

  private Organization createDefaultOrganization(Account defaultAccount) {
    Organization defaultOrg = new Organization();
    defaultOrg.setOrgIdentifier(defaultOrgId);
    defaultOrg.setAccountIdentifier(defaultAccount.getAccountIdentifier());
    defaultOrg.setName(defaultOrgName);
    defaultOrg.setCreatedBy("SYSTEM");
    defaultOrg.setCreatedAt(System.currentTimeMillis());
    defaultOrg.setUpdatedBy("SYSTEM");

    return organizationService.createOrganizationWithRetry(defaultOrg);
  }

  private Project createDefaultProject(Account defaultAccount, Organization defaultOrg) {
    Project defaultProject = new Project();
    defaultProject.setProjectIdentifier(defaultProjectId);
    defaultProject.setOrgIdentifier(defaultOrg.getOrgIdentifier());
    defaultProject.setAccountIdentifier(defaultAccount.getAccountIdentifier());
    defaultProject.setName(defaultProjectName);
    defaultProject.setCreatedBy("SYSTEM");
    defaultProject.setCreatedAt(System.currentTimeMillis());
    defaultProject.setUpdatedBy("SYSTEM");

    return projectService.createProjectWithRetry(defaultProject);
  }

  private UserMapping createUserMapping(
      User createdUser, Account defaultAccount, Organization defaultOrg, Project defaultProject) {
    // Implement user mapping logic
    return userMappingService.createUserMapping(
        createdUser.getUserId(),
        defaultAccount.getAccountIdentifier(),
        defaultOrg.getOrgIdentifier(),
        defaultProject.getProjectIdentifier(),
        true);
  }

  private void fallbackCleanup(
      User createdUser, Account defaultAccount, Organization defaultOrg, Project defaultProject) {
    // Cleanup logic in case of failure
    logger.info("Cleaning up resources due to failure...");

    if (defaultProject != null) {
      projectService.deleteProject(defaultProject.getProjectIdentifier());
    }

    if (defaultOrg != null) {
      organizationService.deleteOrganization(defaultOrg.getOrgIdentifier());
    }

    if (defaultAccount != null) {
      accountService.deleteAccount(defaultAccount.getAccountIdentifier());
    }

    userService.deleteById(createdUser.getId());
  }

  @PostConstruct
  public void startProcess() {
    logger.info("Starting background process for hierarchy service...");
    backgroundProcess();
    logger.info("Background process for hierarchy service finished!");
  }
}
