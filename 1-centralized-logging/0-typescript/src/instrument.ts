/**
 * Sentry instrument — initialize the SDK before NestFactory runs.
 * Sentry must hook into the Node runtime as early as possible to capture all exceptions/traces.
 * Reads env directly because ConfigModule does not exist yet at this point.
 */
import * as Sentry from "@sentry/nestjs"

// Only initialize when SENTRY_DSN is configured (skip if empty — not mandatory for the lab).
const dsn = process.env.SENTRY_DSN
if (dsn) {
    Sentry.init({
        dsn,
        // Trace sample rate — 1.0 = 100% (lab only; production should be < 0.1).
        tracesSampleRate: 1.0,
    })
}
