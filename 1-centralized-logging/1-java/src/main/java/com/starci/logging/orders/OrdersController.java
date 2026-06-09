package com.starci.logging.orders;

import java.util.Map;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/** Orders controller — routes delegate to the service. */
@RestController
@RequestMapping("/api/orders")
public class OrdersController {

    private final OrdersService ordersService;

    public OrdersController(OrdersService ordersService) {
        this.ordersService = ordersService;
    }

    /** GET /api/orders — accepted order; logs INFO to Loki for the query demo. */
    @GetMapping
    public Map<String, String> accept() {
        return ordersService.accept();
    }

    /** GET /api/orders/error — deliberate failure; Sentry captures the exception. */
    @GetMapping("/error")
    public Map<String, String> error() {
        return ordersService.fail();
    }
}
