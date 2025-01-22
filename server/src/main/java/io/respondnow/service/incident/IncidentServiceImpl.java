package io.respondnow.service.incident;

import io.respondnow.dto.incident.CreateRequest;
import io.respondnow.exception.IncidentNotFoundException;
import io.respondnow.exception.InvalidIncidentException;
import io.respondnow.exception.RoleUpdateException;
import io.respondnow.model.incident.*;
import io.respondnow.model.user.UserDetails;
import io.respondnow.repository.IncidentRepository;
import java.time.Instant;
import java.util.*;
import java.util.function.Function;
import java.util.stream.Collectors;
import org.bson.types.ObjectId;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.data.mongodb.core.query.Update;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class IncidentServiceImpl implements IncidentService {

  private static final Logger logger = LoggerFactory.getLogger(IncidentServiceImpl.class);
  @Autowired private IncidentRepository incidentRepository;
  @Autowired private MongoTemplate mongoTemplate;

  @Value("${hierarchy.defaultAccount.id:default_account_id}")
  private String defaultAccountId;

  @Value("${hierarchy.defaultOrg.id:default_org_id}")
  private String defaultOrgId;

  @Value("${hierarchy.defaultProject.id:default_project_id}")
  private String defaultProjectId;

  @NotNull
  private static Timeline getTimeline(
      CreateRequest request, UserDetails currentUser, long createdAt) {

    // Check if channels are available and not empty
    if (request.getChannels() == null || request.getChannels().isEmpty()) {
      throw new IllegalArgumentException("No channels provided in the request.");
    }

    // Get the first channel from the list (you can adjust this logic if multiple channels should be
    // handled differently)
    Channel slackChannel = request.getChannels().get(0);

    // Create a new timeline entry for the Slack channel creation
    Timeline slackChannelTimeline = new Timeline();
    slackChannelTimeline.setType(ChangeType.Slack_Channel_Created);
    slackChannelTimeline.setCreatedAt(createdAt);
    slackChannelTimeline.setUpdatedAt(createdAt);
    slackChannelTimeline.setPreviousState(slackChannel.getId());
    slackChannelTimeline.setCurrentState(
        slackChannel.getId()); // Use the channel ID as the current state

    io.respondnow.model.incident.Slack slack = new io.respondnow.model.incident.Slack();
    slack.setChannelId(slackChannel.getId());
    slack.setChannelName(slackChannel.getName());
    slack.setChannelStatus(slackChannel.getStatus());
    slack.setTeamId(slackChannel.getTeamId());
    slackChannelTimeline.setSlack(slack); // Set the full Slack channel object
    slackChannelTimeline.setUserDetails(
        currentUser); // Set the user details associated with the request
    slackChannelTimeline.setMessage("Slack channel associated with the incident");
    slackChannelTimeline.setAdditionalDetails(null); // Add any additional details if necessary

    return slackChannelTimeline;
  }

  public Incident createIncident(CreateRequest request, UserDetails currentUser) {
    long createdAt = Instant.now().getEpochSecond();
    String incidentId = generateIncidentIdentifier(createdAt);

    // Set default status if not provided
    if (request.getStatus() == null) {
      request.setStatus(Status.Started);
    }

    // Initialize new Incident object
    Incident newIncident = new Incident();
    newIncident.setAccountIdentifier(defaultAccountId);
    newIncident.setOrgIdentifier(defaultOrgId);
    newIncident.setProjectIdentifier(defaultProjectId);
    newIncident.setIdentifier(incidentId);
    newIncident.setName(request.getName());
    newIncident.setDescription(request.getDescription());
    newIncident.setType(request.getType());
    newIncident.setSummary(request.getSummary());
    newIncident.setStatus(request.getStatus());
    newIncident.setSeverity(request.getSeverity());
    newIncident.setIncidentChannel(request.getIncidentChannel());
    newIncident.setChannels(request.getChannels());
    newIncident.setServices(request.getServices());
    newIncident.setFunctionalities(request.getFunctionalities());
    newIncident.setEnvironments(request.getEnvironments());
    newIncident.setAttachments(request.getAttachments());
    newIncident.setCreatedBy(currentUser);
    newIncident.setActive(true);
    newIncident.setCreatedAt(createdAt);
    newIncident.setUpdatedAt(createdAt);
    newIncident.setRoles(request.getRoles());

    // Create the INCIDENT_CREATED timeline entry
    Timeline incidentCreatedTimeline = new Timeline();
    incidentCreatedTimeline.setType(ChangeType.Incident_Created);
    incidentCreatedTimeline.setCreatedAt(createdAt);
    incidentCreatedTimeline.setUpdatedAt(createdAt);
    incidentCreatedTimeline.setPreviousState(request.getStatus().toString());
    incidentCreatedTimeline.setCurrentState(request.getStatus().toString());
    incidentCreatedTimeline.setSlack(
        request.getIncidentChannel() != null ? request.getIncidentChannel().getSlack() : null);
    incidentCreatedTimeline.setUserDetails(currentUser);
    incidentCreatedTimeline.setMessage("Incident created");
    incidentCreatedTimeline.setAdditionalDetails(null); // Add any additional details if necessary
    newIncident.addTimeline(incidentCreatedTimeline);

    // If Incident Channel and Slack Channel details are provided, add a timeline entry
    if (request.getIncidentChannel() != null
        && request.getIncidentChannel().getSlack().getChannelId() != null) {
      Timeline slackChannelTimeline = getTimeline(request, currentUser, createdAt);
      newIncident.addTimeline(slackChannelTimeline);
    }

    // Save and return the new Incident
    return incidentRepository.save(newIncident);
  }

  public Incident updateSummary(String incidentID, String newSummary, UserDetails currentUser)
      throws Exception {
    // Step 1: Retrieve the existing incident by its ID
    Optional<Incident> existingIncident = incidentRepository.findByIdentifier(incidentID);
    if (existingIncident.isEmpty()) {
      throw new Exception("Incident not found with ID: " + incidentID);
    }

    Incident incident = existingIncident.get();
    // Step 2: Get the old summary and prepare the new timeline entry
    String oldSummary = incident.getSummary();

    // Get the current timestamp (in Unix time)
    long ts = Instant.now().getEpochSecond();

    // Update the audit details with the current user and timestamp
    incident.setUpdatedBy(currentUser);
    incident.setUpdatedAt(ts);
    incident.setUpdatedAt(ts);

    // Step 3: Create a new timeline entry for the change
    Timeline timeline = new Timeline();
    timeline.setId(String.valueOf(ts));
    timeline.setType(ChangeType.Summary);
    timeline.setCreatedAt(ts);
    timeline.setUpdatedAt(ts);
    timeline.setUserDetails(currentUser);
    timeline.setPreviousState(oldSummary);
    timeline.setCurrentState(newSummary);

    // Add the timeline entry to the incident's timeline
    incident.getTimelines().add(timeline);

    // Step 4: Update the incident's summary and description
    incident.setSummary(newSummary);
    incident.setDescription(newSummary);

    // Step 5: Update the incident in the database
    Incident updated = updateIncidentById(incident.getId(), incident);
    if (updated == null) {
      throw new Exception("Failed to update incident summary.");
    }

    return updated;
  }

  @Transactional
  public Incident updateIncidentRoles(
      String incidentID, List<Role> newRoleAssignments, UserDetails currentUser)
      throws RoleUpdateException, IncidentNotFoundException {

    logger.info("Starting role update for incident ID: {}", incidentID);

    // Step 1: Retrieve the existing incident by its ID
    Optional<Incident> existingIncidentOpt = incidentRepository.findByIdentifier(incidentID);
    if (existingIncidentOpt.isEmpty()) {
      logger.error("Incident not found with ID: {}", incidentID);
      throw new IncidentNotFoundException("Incident not found with ID: " + incidentID);
    }

    Incident incident = existingIncidentOpt.get();
    List<Role> existingRoles = incident.getRoles();

    // Step 2: Create maps for existing roles
    Map<RoleType, Role> roleTypeToRoleMap =
        existingRoles.stream()
            .filter(role -> role.getRoleType() != null) // Filter out roles with null RoleType
            .collect(
                Collectors.toMap(
                    Role::getRoleType, Function.identity(), (existing, replacement) -> existing));

    Map<String, List<RoleType>> userIdToRoleTypesMap =
        existingRoles.stream()
            .filter(
                role ->
                    role.getUserDetails() != null
                        && role.getUserDetails().getUserId()
                            != null) // Filter out roles with null userDetails or userId
            .collect(
                Collectors.groupingBy(
                    role -> role.getUserDetails().getUserId(),
                    Collectors.mapping(Role::getRoleType, Collectors.toList())));

    // Initialize lists to track changes for the timeline
    List<String> previousStates = new ArrayList<>();
    List<String> currentStates = new ArrayList<>();
    Set<String> affectedUsers = new HashSet<>();

    // Keep a copy of existing roles before modifications for previousState
    List<Role> previousRoleSnapshot = new ArrayList<>(existingRoles);

    // Step 3: Process each new role assignment
    for (Role newRole : newRoleAssignments) {
      RoleType newRoleType = newRole.getRoleType();
      UserDetails newUser = newRole.getUserDetails();
      String newUserId = newUser != null ? newUser.getUserId() : null;

      if (newRoleType == null || newUserId == null) {
        logger.warn("Skipping invalid role assignment: RoleType or UserId is null");
        continue; // Skip if role type or user id is null
      }

      // If the RoleType is already assigned, we need to replace the existing user
      if (roleTypeToRoleMap.containsKey(newRoleType)) {
        Role existingRole = roleTypeToRoleMap.get(newRoleType);
        String existingUserId = existingRole.getUserDetails().getUserId();

        if (!existingUserId.equals(newUserId)) {
          // Remove the existing user for this role
          existingRoles.remove(existingRole);
          previousStates.add(existingRole.toString());
          affectedUsers.add(existingUserId);
          logger.info("Removed role '{}' from user '{}'", newRoleType, existingUserId);
        }
      }

      // If the new role assignment is not already present, assign the new role
      if (!userIdToRoleTypesMap.containsKey(newUserId)
          || !userIdToRoleTypesMap.get(newUserId).contains(newRoleType)) {
        existingRoles.add(newRole);
        currentStates.add(newRole.toString());
        affectedUsers.add(newUserId);
        logger.info("Assigned role '{}' to user '{}'", newRoleType, newUserId);
      } else {
        logger.info(
            "User '{}' already has role '{}' - skipping assignment.", newUserId, newRoleType);
      }
    }

    if (currentStates.isEmpty() && previousStates.isEmpty()) {
      logger.warn("No roles were updated for incident ID: {}", incidentID);
      throw new RoleUpdateException("No roles were updated. Please provide different roles.");
    }

    // Step 4: Get the current timestamp (in Unix time)
    long ts = Instant.now().getEpochSecond();

    // Update the audit details with the current user and timestamp
    incident.setUpdatedBy(currentUser);
    incident.setUpdatedAt(ts);

    // Step 5: Create a new timeline entry for the change
    Timeline timeline = new Timeline();
    timeline.setId(String.valueOf(ts));
    timeline.setType(ChangeType.Roles);
    timeline.setCreatedAt(ts);
    timeline.setUpdatedAt(ts);
    timeline.setUserDetails(currentUser);

    Map<String, Object> roleDetailsMap = new HashMap<>();
    roleDetailsMap.put("previousState", previousRoleSnapshot);
    roleDetailsMap.put("currentState", existingRoles);
    timeline.setAdditionalDetails(roleDetailsMap);

    timeline.setPreviousState(String.join(" | ", previousStates));
    timeline.setCurrentState(String.join(" | ", currentStates));

    logger.info("Updated roles: {}", existingRoles);

    incident.addTimeline(timeline);
    incident.setRoles(existingRoles);
    return incidentRepository.save(incident);
  }

  public Incident updateIncidentSeverity(
      String incidentID, Severity newSeverity, UserDetails currentUser) throws Exception {
    // Step 1: Retrieve the existing incident by its ID
    Optional<Incident> existingIncident = incidentRepository.findByIdentifier(incidentID);
    if (existingIncident.isEmpty()) {
      throw new Exception("Incident not found with ID: " + incidentID);
    }

    Incident incident = existingIncident.get();
    // Step 2: Get the old role and prepare the new timeline entry
    Severity oldSeverity = incident.getSeverity();

    // Get the current timestamp (in Unix time)
    long ts = Instant.now().getEpochSecond();

    // Update the audit details with the current user and timestamp
    incident.setUpdatedBy(currentUser);
    incident.setUpdatedAt(ts);
    incident.setUpdatedAt(ts);

    // Step 3: Create a new timeline entry for the change
    Timeline timeline = new Timeline();
    timeline.setId(String.valueOf(ts));
    timeline.setType(ChangeType.Severity);
    timeline.setCreatedAt(ts);
    timeline.setUpdatedAt(ts);
    timeline.setUserDetails(currentUser);
    timeline.setPreviousState(oldSeverity.toString());
    timeline.setCurrentState(newSeverity.toString());

    // Add the timeline entry to the incident's timeline
    incident.getTimelines().add(timeline);

    // Step 4: Update the incident's summary and description
    incident.setSeverity(newSeverity);

    // Step 5: Update the incident in the database
    Incident updated = updateIncidentById(incident.getId(), incident);
    if (updated == null) {
      throw new Exception("Failed to update incident summary.");
    }

    return updated;
  }

  public Incident addComment(String incidentID, String comment, UserDetails currentUser)
      throws Exception {
    // Step 1: Retrieve the existing incident by its ID
    Optional<Incident> existingIncident = incidentRepository.findByIdentifier(incidentID);
    if (existingIncident.isEmpty()) {
      throw new Exception("Incident not found with ID: " + incidentID);
    }

    Incident incident = existingIncident.get();
    // Step 2: Get the old comment and prepare the new timeline entry
    List<String> comments = incident.getComment();
    if (comments == null) {
      comments = new ArrayList<>();
    }

    // Get the current timestamp (in Unix time)
    long ts = Instant.now().getEpochSecond();

    // Update the audit details with the current user and timestamp
    incident.setUpdatedBy(currentUser);
    incident.setUpdatedAt(ts);
    incident.setUpdatedAt(ts);

    // Step 3: Create a new timeline entry for the change
    Timeline timeline = new Timeline();
    timeline.setId(String.valueOf(ts));
    timeline.setType(ChangeType.Comment);
    timeline.setCreatedAt(ts);
    timeline.setUpdatedAt(ts);
    timeline.setUserDetails(currentUser);
    timeline.setPreviousState(comment);
    timeline.setCurrentState(comment);

    // Add the timeline entry to the incident's timeline
    incident.getTimelines().add(timeline);

    // Step 4: Update the incident's summary and description
    comments.add(comment);
    incident.setComment(comments);

    // Step 5: Update the incident in the database
    Incident updated = updateIncidentById(incident.getId(), incident);
    if (updated == null) {
      throw new Exception("Failed to add a new incident comment.");
    }

    return updated;
  }

  public Incident updateStatus(String incidentID, Status newStatus, UserDetails currentUser)
      throws Exception {
    Optional<Incident> existingIncident = incidentRepository.findByIdentifier(incidentID);
    if (existingIncident.isEmpty()) {
      throw new Exception("Incident not found with ID: " + incidentID);
    }

    Incident incident = existingIncident.get();
    Status oldStatus = incident.getStatus();

    // Get the current timestamp (in Unix time)
    long ts = Instant.now().getEpochSecond();

    // Update the audit details with the current user and timestamp
    incident.setUpdatedBy(currentUser);
    incident.setUpdatedAt(ts);
    incident.setUpdatedAt(ts);

    // Step 3: Create a new timeline entry for the change
    Timeline timeline = new Timeline();
    timeline.setId(String.valueOf(ts));
    timeline.setType(ChangeType.Status);
    timeline.setCreatedAt(ts);
    timeline.setUpdatedAt(ts);
    timeline.setUserDetails(currentUser);
    timeline.setPreviousState(oldStatus.toString());
    timeline.setCurrentState(newStatus.toString());

    // Add the timeline entry to the incident's timeline
    incident.getTimelines().add(timeline);

    // Step 4: Update the incident's summary and description
    incident.setStatus(newStatus);

    // Step 5: Update the incident in the database
    Incident updated = updateIncidentById(incident.getId(), incident);
    if (updated == null) {
      throw new Exception("Failed to update incident summary.");
    }

    return updated;
  }

  public Incident getIncidentById(String id) {
    return incidentRepository
        .findById(String.valueOf(new ObjectId(id)))
        .orElseThrow(() -> new InvalidIncidentException("Incident not found for ID: " + id));
  }

  public Incident getIncidentByIdentifier(String identifier) {
    Criteria criteria = Criteria.where("identifier").is(identifier);
    Query query = new Query(criteria);
    return mongoTemplate.findOne(query, Incident.class);
  }

  public List<Incident> listIncidents(Query query) {
    return mongoTemplate.find(query, Incident.class);
  }

  public long countIncidents(Query query) {
    return mongoTemplate.count(query, Incident.class);
  }

  @Transactional
  public Incident updateIncidentById(String id, Incident incident) {
    validateIncident(incident);
    long now = Instant.now().getEpochSecond();
    incident.setUpdatedAt(now);

    Query query = new Query(Criteria.where("_id").is(new ObjectId(id)));
    Update update =
        new Update()
            .set("name", incident.getName())
            .set("description", incident.getDescription())
            .set("tags", incident.getTags())
            .set("severity", incident.getSeverity())
            .set("status", incident.getStatus())
            .set("active", incident.getActive())
            .set("summary", incident.getSummary())
            .set("comment", incident.getComment())
            .set("services", incident.getServices())
            .set("environments", incident.getEnvironments())
            .set("functionalities", incident.getFunctionalities())
            .set("roles", incident.getRoles())
            .set("stages", incident.getStages())
            .set("timelines", incident.getTimelines())
            .set("channels", incident.getChannels())
            .set("conferenceDetails", incident.getConferenceDetails())
            .set("attachments", incident.getAttachments())
            .set("updatedAt", now)
            .set("updatedBy", incident.getUpdatedBy());

    mongoTemplate.updateFirst(query, update, Incident.class);

    return getIncidentById(id);
  }

  @Transactional
  public void bulkProcessIncidents(List<Incident> createList, List<Incident> updateList) {
    long now = Instant.now().getEpochSecond();

    createList.forEach(
        incident -> {
          incident.setId(null);
          incident.setCreatedAt(now);
          validateIncident(incident);
        });

    updateList.forEach(
        incident -> {
          incident.setUpdatedAt(now);
          validateIncident(incident);

          Query query = new Query(Criteria.where("_id").is(incident.getId()));
          Update update =
              new Update()
                  .set("name", incident.getName())
                  .set("description", incident.getDescription())
                  .set("tags", incident.getTags())
                  .set("severity", incident.getSeverity())
                  .set("status", incident.getStatus())
                  .set("active", incident.getActive())
                  .set("summary", incident.getSummary())
                  .set("comment", incident.getComment())
                  .set("services", incident.getServices())
                  .set("environments", incident.getEnvironments())
                  .set("functionalities", incident.getFunctionalities())
                  .set("roles", incident.getRoles())
                  .set("stages", incident.getStages())
                  .set("timelines", incident.getTimelines())
                  .set("channels", incident.getChannels())
                  .set("conferenceDetails", incident.getConferenceDetails())
                  .set("attachments", incident.getAttachments())
                  .set("updatedAt", now);

          mongoTemplate.updateFirst(query, update, Incident.class);
        });

    incidentRepository.saveAll(createList);
  }

  public void validateIncident(Incident incident) {
    if (incident.getIdentifier() == null || incident.getIdentifier().isEmpty()) {
      throw new InvalidIncidentException("Missing identifier");
    }
    if (incident.getName() == null || incident.getName().isEmpty()) {
      throw new InvalidIncidentException("Missing name");
    }
    if (incident.getAccountIdentifier() == null || incident.getAccountIdentifier().isEmpty()) {
      throw new InvalidIncidentException("Missing account ID");
    }
    if (incident.getType() == null) {
      throw new InvalidIncidentException("Missing incident type");
    }
    if (incident.getStatus() == null) {
      throw new InvalidIncidentException("Missing incident status");
    }
    if (incident.getSeverity() == null) {
      throw new InvalidIncidentException("Missing severity");
    }
    if ((incident.getSummary() == null || incident.getSummary().isEmpty())
        && (incident.getDescription() == null || incident.getDescription().isEmpty())) {
      throw new InvalidIncidentException("Either summary or description must not be empty");
    }
  }

  public List<Type> getIncidentTypes() {
    return List.of(Type.Availability, Type.Latency, Type.Security, Type.Other);
  }

  public List<AttachmentType> getIncidentAttachmentTypes() {
    return List.of(AttachmentType.Link);
  }

  public List<Severity> getIncidentSeverities() {
    return List.of(Severity.SEV0, Severity.SEV1, Severity.SEV2);
  }

  public List<Status> getIncidentStatuses() {
    return List.of(
        Status.Started,
        Status.Acknowledged,
        Status.Investigating,
        Status.Identified,
        Status.Mitigated,
        Status.Resolved);
  }

  public String generateIncidentIdentifier(long createdAt) {
    // Convert createdAt (long) to String and concatenate with a new UUID
    return createdAt + "-" + UUID.randomUUID();
  }
}
