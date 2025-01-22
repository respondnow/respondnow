package io.respondnow.filter;
//
//// import io.respondnow.util.JwtUtil;
// import org.springframework.security.core.context.SecurityContextHolder;
// import org.springframework.web.filter.OncePerRequestFilter;
//
// import javax.servlet.Filter;
// import javax.servlet.FilterChain;
// import javax.servlet.FilterConfig;
// import javax.servlet.ServletException;
// import javax.servlet.annotation.WebFilter;
// import javax.servlet.http.HttpServletRequest;
// import javax.servlet.http.HttpServletResponse;
// import java.io.IOException;
//
// @WebFilter("/api/*")
// public class JwtAuthenticationFilter extends OncePerRequestFilter {
//
//    @Override
//    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response,
// FilterChain filterChain)
//            throws ServletException, IOException {
//
//        String token = request.getHeader("Authorization");
//        if (token != null && token.startsWith("Bearer ")) {
//            String jwt = token.substring(7); // Extract JWT token
//
//            // Validate the token
////            if (JwtUtil.validateToken(jwt, "username")) {
////                // Set user details in security context (can be customized later)
////                SecurityContextHolder.getContext().setAuthentication(new
// UsernamePasswordAuthenticationToken("username", null, null));
////            }
////        }
//
//        filterChain.doFilter(request, response);
//    }
//
////    @Override
////    public void init(FilterConfig filterConfig) throws ServletException {
////    }
////
////    @Override
////    public void destroy() {
////    }
// }
