package io.respondnow.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.web.bind.annotation.*;

@Tag(name = "Metrics", description = "Metrics related operations")
@RestController
@RequestMapping("/metrics")
public class MetricController {

  @Operation(summary = "Get Metrics", description = "This endpoint returns the application metrics")
  @GetMapping
  public String metrics() {
    return "Metrics exposed";
  }
}
