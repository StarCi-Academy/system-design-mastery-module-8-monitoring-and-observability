package com.starci.consulapi;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Entry point for the Spring Boot consul-api proxy.
 */
@SpringBootApplication
public class ConsulApiApplication {

    public static void main(String[] args) {
        SpringApplication.run(ConsulApiApplication.class, args);
    }
}
