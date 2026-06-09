package com.starci.consulapi;

/**
 * JSON body accepted by POST /consul/register.
 */
public record RegisterRequest(String id, String name, String address, int port) {
}
