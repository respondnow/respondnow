package io.respondnow.model.hierarchy;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

// Entity Definitions
@Data
@NoArgsConstructor
@Getter
@Setter
@Document(collection = "accounts")
public class Account {
  @Id private String id;

  @NotBlank private String accountIdentifier;

  @NotBlank private String name;

  @NotNull private Long createdAt;

  private Long updatedAt;

  @NotBlank private String createdBy;

  private String updatedBy;

  private boolean removed;

  // Getters and Setters
}
