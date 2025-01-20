package io.respondnow.controller;

import io.respondnow.dto.auth.*;
import io.respondnow.exception.EmailAlreadyExistsException;
import io.respondnow.exception.UserNotFoundException;
import io.respondnow.model.user.User;
import io.respondnow.service.auth.AuthService;
import io.respondnow.service.hierarchy.UserMappingService;
import io.respondnow.util.JWTUtil;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import jakarta.validation.Valid;
import java.util.UUID;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/auth")
public class AuthController {

  @Autowired private AuthService authService;

  @Autowired private UserMappingService userMappingService;

  @Autowired private JWTUtil jwtUtil;

  private static final Logger logger = LoggerFactory.getLogger(AuthController.class);

  @Operation(summary = "Sign up a new user")
  @ApiResponses({
    @ApiResponse(responseCode = "201", description = "User signed up successfully"),
    @ApiResponse(responseCode = "400", description = "Bad Request"),
    @ApiResponse(responseCode = "409", description = "Conflict - Email already exists")
  })
  @PostMapping("/signup")
  public ResponseEntity<SignupResponseDTO> signup(@RequestBody @Valid AddUserInput input) {
    try {
      User user = authService.signup(input);
      String token = jwtUtil.generateToken(user.getName(), user.getUserId(), user.getEmail());

      SignupResponseDTO response =
          new SignupResponseDTO("success", "User registered successfully", token, user);
      return ResponseEntity.status(201).body(response);
    } catch (EmailAlreadyExistsException e) {
      // Handle the case where email already exists
      SignupResponseDTO response =
          new SignupResponseDTO("error", "Email already exists", null, null);
      return ResponseEntity.status(409).body(response); // 409 Conflict
    } catch (Exception e) {
      // Handle any other errors (e.g. bad request, unexpected errors)
      SignupResponseDTO response = new SignupResponseDTO("error", "Bad Request", null, null);
      return ResponseEntity.status(400).body(response); // 400 Bad Request
    }
  }

  @Operation(summary = "User login")
  @ApiResponses({
    @ApiResponse(responseCode = "200", description = "Login successful"),
    @ApiResponse(responseCode = "400", description = "Bad Request"),
    @ApiResponse(responseCode = "401", description = "Unauthorized - Invalid credentials")
  })
  @PostMapping("/login")
  public ResponseEntity<LoginResponseDTO> login(@RequestBody @Valid LoginUserInput input) {
    try {
      User user = authService.login(input);
      String token = jwtUtil.generateToken(user.getName(), user.getUserId(), user.getEmail());

      if (user.getChangePasswordRequired()) {
        LoginResponseData data = new LoginResponseData(token, user.getLastLoginAt(), true);
        LoginResponseDTO response =
            new LoginResponseDTO("success", "Change Password is required", data);
        return ResponseEntity.ok(response);
      }

      LoginResponseData data =
          new LoginResponseData(
              token, System.currentTimeMillis(), user.getChangePasswordRequired());
      LoginResponseDTO response = new LoginResponseDTO("success", "Login successful", data);
      return ResponseEntity.ok(response);
    } catch (UserNotFoundException e) {
      LoginResponseDTO response = new LoginResponseDTO("error", "User not found", null);
      return ResponseEntity.status(404).body(response); // 404 Not Found
    } catch (Exception e) {
      logger.info(e.getMessage());
      LoginResponseDTO response = new LoginResponseDTO("error", "Bad Request", null);
      return ResponseEntity.status(400).body(response); // 400 Bad Request
    }
  }

  @Operation(summary = "Change user password")
  @ApiResponses({
    @ApiResponse(responseCode = "200", description = "Password changed successfully"),
    @ApiResponse(responseCode = "400", description = "Bad Request"),
    @ApiResponse(responseCode = "401", description = "Unauthorized - Invalid credentials"),
    @ApiResponse(responseCode = "404", description = "Not Found - User not found")
  })
  @PostMapping("/changePassword")
  public ResponseEntity<ChangePasswordResponseDTO> changePassword(
      @RequestBody @Valid ChangePasswordInput input) {
    try {
      User user = authService.changePassword(input);
      String token = jwtUtil.generateToken(user.getName(), user.getUserId(), user.getEmail());

      ChangePasswordResponseData data =
          new ChangePasswordResponseData(token, System.currentTimeMillis());
      ChangePasswordResponseDTO response =
          new ChangePasswordResponseDTO("success", "Password changed successfully", data);
      return ResponseEntity.ok(response);
    } catch (UserNotFoundException e) {
      ChangePasswordResponseDTO response =
          new ChangePasswordResponseDTO("error", "User not found", null);
      return ResponseEntity.status(404).body(response); // 404 Not Found
    } catch (Exception e) {
      ChangePasswordResponseDTO response =
          new ChangePasswordResponseDTO("error", "Bad Request", null);
      return ResponseEntity.status(400).body(response); // 404 Not Found
    }
  }

  @Operation(summary = "Get User Mappings")
  @ApiResponses(
      value = {
        @ApiResponse(
            responseCode = "200",
            description = "User mappings retrieved successfully",
            content = @Content(schema = @Schema(implementation = GetUserMappingResponseDTO.class))),
        @ApiResponse(
            responseCode = "400",
            description = "Invalid input",
            content = @Content(schema = @Schema(implementation = GetUserMappingResponseDTO.class))),
        @ApiResponse(
            responseCode = "500",
            description = "Internal Server Error",
            content = @Content(schema = @Schema(implementation = GetUserMappingResponseDTO.class)))
      })
  @GetMapping("/userMapping")
  public ResponseEntity<GetUserMappingResponseDTO> getUserMappings(
      @RequestParam(value = "correlationId", required = false) String correlationId,
      @RequestParam(value = "userId", required = true) String userId) {

    // Generate correlationId if not provided
    if (correlationId == null || correlationId.isEmpty()) {
      correlationId = UUID.randomUUID().toString();
    }
    GetUserMappingResponseDTO response = new GetUserMappingResponseDTO();
    response.setCorrelationId(correlationId);

    if (userId == null || userId.isEmpty()) {
      response.setMessage("userId is required in the query");
      response.setStatus("ERROR");
      return ResponseEntity.badRequest().body(response);
    }
    try {
      UserMappingData userMappingData = userMappingService.getUserMappings(correlationId, userId);
      response.setData(userMappingData);
      response.setMessage("User mappings retrieved successfully");
      response.setStatus("SUCCESS");
      return ResponseEntity.ok(response);
    } catch (RuntimeException e) {
      response.setMessage("Error: " + e.getMessage());
      response.setStatus("ERROR");
      return ResponseEntity.internalServerError().body(response);
    }
  }
}
