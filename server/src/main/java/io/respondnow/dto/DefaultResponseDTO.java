package io.respondnow.dto;

import lombok.*;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class DefaultResponseDTO {
  private String correlationId;
  private String message;
}
