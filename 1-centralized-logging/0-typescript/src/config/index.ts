/**
 * Config namespace factories — read process.env only inside the factory.
 */
import { registerAs } from "@nestjs/config"

export interface AppConfig {
    port: number
}

export interface LoggingConfig {
    lokiUrl: string
    appName: string
    env: string
}

// App runtime config — `app` namespace.
export const appConfig = registerAs(
    "app",
    (): AppConfig => ({
        port: Number(process.env.PORT) || 3000,
    }),
)

// Logging transport config (Loki labels + endpoint) — `logging` namespace.
export const loggingConfig = registerAs(
    "logging",
    (): LoggingConfig => ({
        lokiUrl: process.env.LOKI_URL ?? "http://loki:3100",
        appName: process.env.LOG_APP_NAME ?? "logging-api",
        env: process.env.LOG_ENV ?? "lab",
    }),
)
