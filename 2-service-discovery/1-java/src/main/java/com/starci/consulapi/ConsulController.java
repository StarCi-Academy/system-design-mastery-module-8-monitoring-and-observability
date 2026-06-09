package com.starci.consulapi;

import java.util.Map;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestClient;

/**
 * Consul proxy controller — register / health / deregister via the Consul Agent HTTP API.
 */
@RestController
@RequestMapping("/consul")
public class ConsulController {

    private final RestClient restClient;
    private final String baseUrl;

    public ConsulController(@Value("${consul.base-url:http://localhost:8500}") String baseUrl) {
        this.baseUrl = baseUrl;
        this.restClient = RestClient.create();
    }

    /**
     * Register a service instance into the Consul catalog via the Agent API.
     */
    @PostMapping("/register")
    @ResponseStatus(HttpStatus.CREATED)
    public Map<String, String> register(@RequestBody RegisterRequest body) {
        // Consul requires PascalCase fields: ID, Name, Address, Port.
        Map<String, Object> payload = Map.of(
                "ID", body.id(),
                "Name", body.name(),
                "Address", body.address(),
                "Port", body.port());
        restClient.put()
                .uri(baseUrl + "/v1/agent/service/register")
                .contentType(MediaType.APPLICATION_JSON)
                .body(payload)
                .retrieve()
                .toBodilessEntity();
        return Map.of("status", "registered");
    }

    /**
     * List passing instances of a service by logical name.
     */
    @GetMapping("/health/{service}")
    public Object health(@PathVariable String service) {
        // Query by logical service name; the health endpoint returns only passing instances.
        return restClient.get()
                .uri(baseUrl + "/v1/health/service/" + service)
                .retrieve()
                .body(Object.class);
    }

    /**
     * Deregister a service instance from the Agent by its exact instance ID.
     */
    @PostMapping("/deregister/{id}")
    @ResponseStatus(HttpStatus.CREATED)
    public Map<String, String> deregister(@PathVariable String id) {
        // Remove exactly this instance from the catalog by its instance ID.
        restClient.put()
                .uri(baseUrl + "/v1/agent/service/deregister/" + id)
                .retrieve()
                .toBodilessEntity();
        return Map.of("status", "deregistered");
    }
}
