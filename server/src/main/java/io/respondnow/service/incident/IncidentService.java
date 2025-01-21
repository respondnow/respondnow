package io.respondnow.service.incident;

import io.respondnow.dto.incident.CreateRequest;
import io.respondnow.model.incident.*;
import io.respondnow.model.user.UserDetails;
import java.util.List;
import org.springframework.data.mongodb.core.query.Query;

public interface IncidentService {
  Incident createIncident(CreateRequest request, UserDetails currentUser);

  Incident getIncidentById(String id);
  Incident getIncidentByIdentifier(String id);

  List<Incident> listIncidents(Query query);

  long countIncidents(Query query);

  Incident updateIncidentById(String id, Incident incident);

  void bulkProcessIncidents(List<Incident> createList, List<Incident> updateList);

  void validateIncident(Incident incident);

  List<Type> getIncidentTypes();

  List<AttachmentType> getIncidentAttachmentTypes();

  List<Severity> getIncidentSeverities();

  List<Status> getIncidentStatuses();

  String generateIncidentIdentifier(long createdAt);

  Incident updateSummary(String incidentID, String newSummary, UserDetails currentUser)
      throws Exception;

  Incident updateStatus(String incidentID, Status newStatus, UserDetails currentUser)
          throws Exception;

  Incident addComment(String incidentID, String comment, UserDetails currentUser) throws Exception;
}
