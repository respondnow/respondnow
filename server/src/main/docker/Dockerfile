FROM openjdk:11-jre-slim

# Add the jar file to the container
COPY target/respondnow-0.0.1-SNAPSHOT.jar /app/respondnow.jar

# Expose the application port
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["java", "-jar", "/app/respondnow.jar"]
