package io.respondnow.service.incident;

import io.respondnow.model.incident.*;
import java.util.List;
import org.springframework.data.mongodb.core.query.Query;

public interface IncidentService {
  Incident createIncident(Incident incident);

  Incident getIncidentById(String id);

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
}
