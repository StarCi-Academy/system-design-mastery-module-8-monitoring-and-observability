package com.starci.tracing;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Spring Boot entry point for the distributed-tracing lab.
 */
@SpringBootApplication
public class TracingApiApplication {

    public static void main(String[] args) {
        SpringApplication.run(TracingApiApplication.class, args);
    }
}
