package io.respondnow.config;

import io.respondnow.security.JWTAuthenticationFilter;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

@Configuration
@EnableWebSecurity
public class SecurityConfig extends WebSecurityConfigurerAdapter {

  private final JWTAuthenticationFilter jwtAuthenticationFilter;

  public SecurityConfig(JWTAuthenticationFilter jwtAuthenticationFilter) {
    this.jwtAuthenticationFilter = jwtAuthenticationFilter;
  }

  @Override
  protected void configure(HttpSecurity http) throws Exception {
    http.csrf()
        .disable() // Disable CSRF for REST API
        .authorizeRequests()
        .antMatchers(
            "/public/**",
            "/login",
            "/auth/login",
            "/signup",
            "/auth/signup",
            "/status",
            "/auth/changePassword")
        .permitAll() // Allow public access to specific endpoints
        // Allow Swagger UI and API docs endpoint to be accessed without authentication
        .antMatchers("/swagger-ui/**", "/v3/api-docs/**", "/swagger-ui.html")
        .permitAll()
        .anyRequest()
        .authenticated() // All other requests require authentication
        .and()
        .addFilterBefore(
            jwtAuthenticationFilter, UsernamePasswordAuthenticationFilter.class) // Add JWT filter

        // Disable the default login page and form login
        .formLogin()
        .disable()
        .httpBasic()
        .disable();
  }
}
