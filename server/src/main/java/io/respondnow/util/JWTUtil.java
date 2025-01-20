package io.respondnow.util;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.security.Keys;
import java.nio.charset.StandardCharsets;
import java.security.Key;
import java.util.Date;
import javax.crypto.SecretKey;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class JWTUtil {

  @Value("${jwt.secret}")
  private String secretKey;

  @Value("${jwt.expiration:86400}")
  private long expirationTime;

  private Key getSigningKey() {
    return Keys.hmacShaKeyFor(secretKey.getBytes());
  }

  // Generate JWT Token
  public String generateToken(String username, String userId, String email) {
    return Jwts.builder()
        .setSubject(username)
        .claim("username", userId)
        .claim("email", email)
        .claim("name", username)
        .setIssuedAt(new Date())
        .setExpiration(
            new Date(System.currentTimeMillis() + expirationTime * 1000)) // Convert to milliseconds
        .signWith(getSigningKey(), SignatureAlgorithm.HS256)
        .compact();
  }

  // Validate JWT Token
  public boolean validateToken(String token, String username) {
    String usernameFromToken = getUsernameFromToken(token);
    return (usernameFromToken.equals(username) && !isTokenExpired(token));
  }

  // Extract username from JWT token
  public String getUsernameFromToken(String token) {
    Claims claims = getClaimsFromToken(token);
    return claims.getSubject();
  }

  // Check if the token is expired
  public boolean isTokenExpired(String token) {
    Claims claims = getClaimsFromToken(token);
    return claims.getExpiration().before(new Date());
  }

  // Extract claims from JWT token
  private Claims getClaimsFromToken(String token) {
    byte[] keyBytes = secretKey.getBytes(StandardCharsets.UTF_8); // Convert the secret to bytes
    SecretKey key = Keys.hmacShaKeyFor(keyBytes);

    return Jwts.parser().setSigningKey(key).build().parseClaimsJws(token).getBody();
  }
}
