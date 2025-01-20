package io.respondnow.security;

import io.respondnow.util.JWTUtil;
import java.io.IOException;
import java.util.ArrayList;
import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

@Component
public class JWTAuthenticationFilter extends OncePerRequestFilter {

  @Autowired private JWTUtil jwtUtil;

  @Override
  protected void doFilterInternal(
      HttpServletRequest request,
      javax.servlet.http.HttpServletResponse response,
      FilterChain filterChain)
      throws ServletException, IOException {

    String token = getJWTFromRequest(request);
    if (token != null && jwtUtil.validateToken(token, jwtUtil.getUsernameFromToken(token))) {
      UsernamePasswordAuthenticationToken authentication =
          new UsernamePasswordAuthenticationToken(
              jwtUtil.getUsernameFromToken(token), null, new ArrayList<>() // No authorities needed
              );
      authentication.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
      SecurityContextHolder.getContext().setAuthentication(authentication);
    }

    filterChain.doFilter(request, response);
  }

  // Extract JWT from the request
  private String getJWTFromRequest(HttpServletRequest request) {
    String bearerToken = request.getHeader("Authorization");
    if (bearerToken != null && bearerToken.startsWith("Bearer ")) {
      return bearerToken.substring(7);
    }
    return null;
  }
}
