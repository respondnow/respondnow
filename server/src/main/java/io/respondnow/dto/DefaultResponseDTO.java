package io.respondnow.dto;

import lombok.*;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class DefaultResponseDTO {
  private String correlationId;
  private String message;
}
