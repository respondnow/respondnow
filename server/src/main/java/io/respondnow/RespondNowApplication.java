package io.respondnow;

import io.swagger.v3.oas.annotations.OpenAPIDefinition;
import io.swagger.v3.oas.annotations.info.Info;
import io.swagger.v3.oas.annotations.info.License;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

@OpenAPIDefinition(
    info =
        @Info(
            title = "RespondNow API",
            version = "0.0.1",
            description =
                "API for the RespondNow Application that provides incident management, user authentication, and more.",
            license =
                @License(name = "Apache-2.0 license", url = "http://www.apache.org/licenses/")))
@SpringBootApplication
@EnableAsync
public class RespondNowApplication {

  public static void main(String[] args) {
    SpringApplication.run(RespondNowApplication.class, args);
  }
}
