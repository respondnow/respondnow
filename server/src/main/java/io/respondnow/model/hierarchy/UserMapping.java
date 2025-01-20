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
@Document(collection = "userMappings")
public class UserMapping {
  @Id private String id;

  @NotBlank private String userId;

  @NotBlank private String accountIdentifier;

  private String orgIdentifier;

  private String projectIdentifier;

  @NotNull private Long createdAt;

  private Long updatedAt;

  private boolean removed;

  private boolean isDefault;
}
