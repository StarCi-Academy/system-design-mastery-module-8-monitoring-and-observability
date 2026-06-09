FROM maven:3.9-eclipse-temurin-21 AS builder
WORKDIR /build
COPY pom.xml ./
# Pre-fetch dependencies for faster, cache-friendly rebuilds.
RUN mvn -B -q dependency:go-offline
COPY src ./src
RUN mvn -B -q clean package -DskipTests

FROM eclipse-temurin:21-jre
WORKDIR /app
# OpenTelemetry Java agent provides zero-code auto-instrumentation for the
# inbound Spring WebMVC HTTP span; the manual child span is added in code.
ADD https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v2.6.0/opentelemetry-javaagent.jar /app/opentelemetry-javaagent.jar
COPY --from=builder /build/target/tracing-api-0.0.1.jar /app/app.jar
# Attach javaagent and cap heap at 256 MB so the container fits a dev-machine memory budget.
ENV JAVA_TOOL_OPTIONS="-javaagent:/app/opentelemetry-javaagent.jar -Xmx256m -Xms128m"
EXPOSE 3000
ENTRYPOINT ["java", "-jar", "/app/app.jar"]
