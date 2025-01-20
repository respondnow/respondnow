package io.respondnow.service.hierarchy;

import io.respondnow.dto.auth.UserIdentifiers;
import io.respondnow.dto.auth.UserMappingData;
import io.respondnow.exception.UserMappingNotFoundException;
import io.respondnow.model.hierarchy.Account;
import io.respondnow.model.hierarchy.Organization;
import io.respondnow.model.hierarchy.Project;
import io.respondnow.model.hierarchy.UserMapping;
import io.respondnow.repository.UserMappingRepository;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class UserMappingServiceImpl implements UserMappingService {
  @Autowired private UserMappingRepository userMappingRepository;

  @Autowired private AccountService accountService;
  @Autowired private OrganizationService organizationService;
  @Autowired private ProjectService projectService;

  public List<UserMapping> findAll() {
    return userMappingRepository.findAll();
  }

  public Optional<UserMapping> findById(String id) {
    return userMappingRepository.findById(id);
  }

  public UserMapping save(UserMapping userMapping) {
    return userMappingRepository.save(userMapping);
  }

  public void deleteById(String id) {
    userMappingRepository.deleteById(id);
  }

  @Transactional
  public UserMapping createUserMapping(
      String userId,
      String accountIdentifier,
      String orgIdentifier,
      String projectIdentifier,
      boolean isDefault) {

    UserMapping userMap = new UserMapping();
    userMap.setUserId(userId);
    userMap.setAccountIdentifier(accountIdentifier);
    userMap.setDefault(isDefault);
    userMap.setCreatedAt(System.currentTimeMillis());

    if (!orgIdentifier.isEmpty()) {
      userMap.setOrgIdentifier(orgIdentifier);
    }

    if (!projectIdentifier.isEmpty()) {
      userMap.setProjectIdentifier(projectIdentifier);
    }

    return userMappingRepository.save(userMap);
  }

  public UserMappingData getUserMappings(String correlationId, String userId) {
    // Logic to fetch mappings
    List<UserMapping> mappings = userMappingRepository.findByUserId(userId);
    if (mappings.isEmpty()) {
      throw new UserMappingNotFoundException("No mappings found for the user");
    }

    UserIdentifiers defaultMapping = null;
    List<UserIdentifiers> allMappings = new ArrayList<>();

    for (UserMapping userMapping : mappings) {
      Account account = accountService.getAccount(userMapping.getAccountIdentifier());
      Organization org = organizationService.getOrganization(userMapping.getOrgIdentifier());
      Project project = projectService.getProject(userMapping.getProjectIdentifier());

      UserIdentifiers mappingIdentifiers = new UserIdentifiers();
      mappingIdentifiers.setAccountIdentifier(userMapping.getAccountIdentifier());
      mappingIdentifiers.setAccountName(account.getName());
      mappingIdentifiers.setOrgIdentifier(userMapping.getOrgIdentifier());
      mappingIdentifiers.setOrgName(org.getName());
      mappingIdentifiers.setProjectIdentifier(userMapping.getProjectIdentifier());
      mappingIdentifiers.setProjectName(project.getName());

      allMappings.add(mappingIdentifiers);

      if (userMapping.isDefault()) {
        defaultMapping = mappingIdentifiers;
      }
    }

    UserMappingData userMappingData = new UserMappingData();
    userMappingData.setDefaultMapping(defaultMapping);
    userMappingData.setMappings(allMappings);

    return userMappingData;
  }
}
