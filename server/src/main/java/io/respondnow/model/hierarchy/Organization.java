package io.respondnow.model.hierarchy;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import javax.persistence.Id;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.mongodb.core.mapping.Document;

@Data
@NoArgsConstructor
@Getter
@Setter
@Document(collection = "organizations")
public class Organization {
  @Id private String id;

  @NotBlank private String orgIdentifier;

  @NotBlank private String name;

  @NotBlank private String accountIdentifier;

  @NotNull private Long createdAt;

  private Long updatedAt;

  @NotBlank private String createdBy;

  private String updatedBy;

  private boolean removed;

  // Getters and Setters
}
