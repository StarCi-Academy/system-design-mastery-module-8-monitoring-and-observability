# Multi-stage build for the Spring Boot consul-api proxy.
FROM maven:3.9-eclipse-temurin-21 AS builder
WORKDIR /app
COPY pom.xml ./
RUN mvn -q dependency:go-offline
COPY src ./src
RUN mvn -q clean package -DskipTests

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /app/target/consul-api-0.0.1.jar ./app.jar
EXPOSE 3000
CMD ["java", "-jar", "app.jar"]
