package io.respondnow.model.hierarchy;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

@Data
@NoArgsConstructor
@Getter
@Setter
@Document(collection = "projects")
public class Project {
  @Id private String id;

  @NotBlank private String projectIdentifier;

  @NotBlank private String name;

  @NotBlank private String orgIdentifier;

  @NotBlank private String accountIdentifier;

  @NotNull private Long createdAt;

  private Long updatedAt;

  @NotBlank private String createdBy;

  private String updatedBy;

  private boolean removed;

  // Getters and Setters
}
