package academy.starci.metricsapi;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Spring Boot entry point for the Prometheus metrics lesson (Micrometer + Actuator).
 */
@SpringBootApplication
public class MetricsApiApplication {

    public static void main(String[] args) {
        SpringApplication.run(MetricsApiApplication.class, args);
    }
}
