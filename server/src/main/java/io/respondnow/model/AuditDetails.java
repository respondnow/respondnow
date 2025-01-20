package io.respondnow.model;

import io.respondnow.model.user.User;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Data
@NoArgsConstructor
@Getter
@Setter
public class AuditDetails {

  private Long createdAt;
  private Long updatedAt;
  private User createdBy;
  private User updatedBy;
  private Long removedAt;
  private Boolean removed;
}
