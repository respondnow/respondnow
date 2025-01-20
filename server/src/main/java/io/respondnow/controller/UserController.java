package io.respondnow.controller;

import io.respondnow.model.user.User;
import io.respondnow.service.user.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import java.util.List;
import java.util.Optional;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@Tag(name = "Users", description = "User-related operations")
@RestController
@RequestMapping("/users")
public class UserController {

  @Autowired private UserService userService;

  // Create a new user
  @Operation(summary = "Create a new user", description = "This endpoint creates a new user")
  @PostMapping("/create")
  public User createUser(@RequestBody User user) {
    return userService.save(user);
  }

  // Get all users
  @Operation(summary = "Get all users", description = "This endpoint retrieves a list of all users")
  @GetMapping("/list")
  public List<User> getAllUsers() {
    return userService.findAll();
  }

  // Get user by ID
  @Operation(summary = "Get user by ID", description = "This endpoint retrieves a user by their ID")
  @GetMapping("/{id}")
  public Optional<User> getUserById(@PathVariable String id) {
    return userService.findById(id);
  }
}
