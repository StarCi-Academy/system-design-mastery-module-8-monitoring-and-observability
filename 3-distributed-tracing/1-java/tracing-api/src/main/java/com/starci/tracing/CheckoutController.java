package com.starci.tracing;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * HTTP controller — the inbound span is auto-instrumented; the handler only adds the child span.
 */
@RestController
@RequestMapping("/checkout")
public class CheckoutController {

    private final CheckoutService checkout;

    public CheckoutController(CheckoutService checkout) {
        this.checkout = checkout;
    }

    // Logic: inbound HTTP span is auto-instrumented; the handler only adds the child span.
    @GetMapping
    public CheckoutResult checkout() {
        return checkout.simulateCheckout();
    }
}
