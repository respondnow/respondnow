package io.respondnow.model.user;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;
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
@Document(collection = "users")
public class User {

  @Id private String id;

  @NotBlank private String name;

  @NotBlank private String userId;

  @Email private String email;

  @NotBlank private String password;

  private Boolean active;

  private Boolean changePasswordRequired;

  private Long createdAt;

  private Long lastLoginAt;

  private Long removedAt;

  private Boolean removed;

  private String createdBy;

  private String updatedBy;

  private Long updatedAt;
}
