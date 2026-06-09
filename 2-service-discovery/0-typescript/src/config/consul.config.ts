/**
 * Consul HTTP API endpoint config — `consul` namespace.
 */
import {
    registerAs,
} from "@nestjs/config"

export interface ConsulConfig {
    baseUrl: string
}

export const consulConfig = registerAs(
    "consul",
    (): ConsulConfig => ({
        baseUrl: process.env.CONSUL_BASE_URL ?? "http://localhost:8500",
    }),
)
