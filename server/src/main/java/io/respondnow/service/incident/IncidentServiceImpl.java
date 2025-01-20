package io.respondnow.service.incident;

import io.respondnow.exception.InvalidIncidentException;
import io.respondnow.model.incident.*;
import io.respondnow.repository.IncidentRepository;
import java.time.Instant;
import java.util.List;
import java.util.UUID;
import org.bson.types.ObjectId;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.data.mongodb.core.query.Update;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class IncidentServiceImpl implements IncidentService {

  @Autowired private IncidentRepository incidentRepository;
  @Autowired private MongoTemplate mongoTemplate;

  //  private final String accountIdentifier;
  //  private final String orgIdentifier;
  //  private final String projectIdentifier;
  //
  //  public IncidentServiceImpl(
  //      String accountIdentifier, String orgIdentifier, String projectIdentifier) {
  //    this.accountIdentifier = accountIdentifier;
  //    this.orgIdentifier = orgIdentifier;
  //    this.projectIdentifier = projectIdentifier;
  //  }

  @Transactional
  public Incident createIncident(Incident incident) {
    validateIncident(incident);
    incident.setId(null);
    long now = Instant.now().getEpochSecond();
    incident.setCreatedAt(now);
    incident.setUpdatedAt(now);
    return incidentRepository.save(incident);
  }

  //  public Incident createIncident(IncidentCreateRequest request, UserDetails currentUser) {
  //    long createdAt = Instant.now().getEpochSecond();
  //    String incidentId = generateIncidentIdentifier(createdAt);
  //
  //    Incident newIncident = new Incident();
  //    newIncident.setIdentifier(incidentId);
  //    newIncident.setSummary(request.getSummary());
  //    newIncident.setStatus("Started");
  //    newIncident.setSeverity(request.getSeverity());
  //    newIncident.setIncidentChannel(request.getIncidentChannel());
  //    newIncident.setCreatedBy(currentUser);
  //    newIncident.setCreatedAt(createdAt);
  //    newIncident.setUpdatedAt(createdAt);
  //
  //    // Create the timeline for the incident
  //    Timeline incidentCreatedTimeline =
  //        new Timeline("IncidentCreated", createdAt, null, null, currentUser);
  //    newIncident.addTimeline(incidentCreatedTimeline);
  //
  //    // If Slack channel is provided, handle it here
  //    if (request.getIncidentChannel() != null
  //        && request.getIncidentChannel().getSlackChannelId() != null) {
  //      // Assuming `slackClient.createChannel()` will return a channel object with the Slack
  // channel
  //      // ID
  //      String slackChannelId =
  //          slackClient.createChannel(request.getIncidentChannel().getSlackChannelId());
  //      newIncident.setIncidentChannel(new IncidentChannel(slackChannelId));
  //      Timeline slackCreatedTimeline =
  //          new Timeline("SlackChannelCreated", createdAt, null, slackChannelId, currentUser);
  //      newIncident.addTimeline(slackCreatedTimeline);
  //    }
  //
  //    return incidentRepository.save(newIncident);
  //  }

  public Incident getIncidentById(String id) {
    return incidentRepository
        .findById(String.valueOf(new ObjectId(id)))
        .orElseThrow(() -> new InvalidIncidentException("Incident not found for ID: " + id));
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
            .set("updatedAt", now);

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
    return List.of(Type.AVAILABILITY, Type.LATENCY, Type.SECURITY, Type.OTHER);
  }

  public List<AttachmentType> getIncidentAttachmentTypes() {
    return List.of(AttachmentType.LINK);
  }

  public List<Severity> getIncidentSeverities() {
    return List.of(Severity.SEV0, Severity.SEV1, Severity.SEV2);
  }

  public List<Status> getIncidentStatuses() {
    return List.of(
        Status.STARTED,
        Status.ACKNOWLEDGED,
        Status.INVESTIGATING,
        Status.IDENTIFIED,
        Status.MITIGATED,
        Status.RESOLVED);
  }

  public String generateIncidentIdentifier(long createdAt) {
    // Convert createdAt (long) to String and concatenate with a new UUID
    return createdAt + "-" + UUID.randomUUID();
  }
}
