# Use multi-arch compatible Maven image to build the application
FROM maven:3.9.6-eclipse-temurin-11 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the pom.xml and download dependencies (to leverage Docker cache)
COPY pom.xml .
RUN mvn dependency:go-offline

# Copy the rest of the project files
COPY . .

# Build the application
RUN mvn clean package -DskipTests

# Use a multi-arch compatible OpenJDK 11 JRE Slim image
FROM openjdk:11-jre-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the built jar file from the build stage
COPY --from=build /app/target/*.jar app.jar

# Expose the application port
EXPOSE 8080

# Run the application
ENTRYPOINT ["java", "-jar", "app.jar"]
