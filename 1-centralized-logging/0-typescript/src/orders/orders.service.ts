/**
 * Orders service — generates a trace and writes demo logs.
 */
import { Injectable, Logger } from "@nestjs/common"

@Injectable()
export class OrdersService {
    private readonly logger = new Logger(OrdersService.name)

    /**
     * Simulate a successful order; log INFO with a trace for Loki lookup.
     */
    accept(): { status: string; trace: string } {
        const trace = `trace-${Date.now()}`
        // HOSTNAME in containers helps distinguish replicas in demos.
        this.logger.log(`[${process.env.HOSTNAME ?? "local"}] Order accepted trace=${trace}`)
        return { status: "accepted", trace }
    }

    /**
     * Simulate an order failure; log ERROR then throw a real Error.
     * SentryGlobalFilter catches the exception and reports it to Sentry.
     */
    fail(): never {
        this.logger.error(`[${process.env.HOSTNAME ?? "local"}] Simulated order failure for Loki + Sentry demo`)
        // Throw a real Error — SentryGlobalFilter catches this and reports to the Sentry dashboard.
        throw new Error("Simulated order processing failure")
    }
}
