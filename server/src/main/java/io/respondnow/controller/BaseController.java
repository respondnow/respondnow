package io.respondnow.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.web.bind.annotation.*;

@Tag(name = "Base Operations", description = "Base routes such as status and version information")
@RestController
@RequestMapping("/")
public class BaseController {
  @Operation(
      summary = "Check Server Status",
      description = "This endpoint provides the status of the server")
  @GetMapping("/status")
  public String statusHandler() {
    return "Server is up and running";
  }
}
