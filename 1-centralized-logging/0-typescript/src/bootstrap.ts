/**
 * Bootstrap the Nest app — Winston (Console + Loki transports), listen on the port.
 */
import { Logger } from "@nestjs/common"
import { NestFactory } from "@nestjs/core"
import { WinstonModule } from "nest-winston"
import * as winston from "winston"
import LokiTransport from "winston-loki"
import { AppModule } from "./app.module"

export async function bootstrap(): Promise<void> {
    const port = Number(process.env.PORT) || 3000
    const lokiUrl = process.env.LOKI_URL ?? "http://loki:3100"
    const appName = process.env.LOG_APP_NAME ?? "logging-api"
    const env = process.env.LOG_ENV ?? "lab"

    // Winston replaces the default NestJS logger: Console for dev, Loki for centralized logging.
    const app = await NestFactory.create(AppModule, {
        logger: WinstonModule.createLogger({
            transports: [
                // Console transport — prints logs to stdout (Docker / dev terminal).
                new winston.transports.Console({
                    format: winston.format.combine(
                        winston.format.timestamp(),
                        winston.format.json(),
                    ),
                }),
                // Loki transport — pushes logs via the HTTP push API to Loki.
                new LokiTransport({
                    host: lokiUrl,
                    labels: { app: appName, env },
                    json: true,
                    // Surface transport-level connection errors on stdout instead of swallowing them.
                    onConnectionError: (err): void => console.error("loki transport error", err),
                }),
            ],
        }),
    })

    // Bind 0.0.0.0 for Docker.
    await app.listen(port, "0.0.0.0")
    Logger.log(`logging-api listening on :${port} | loki=${lokiUrl}`, "Bootstrap")
}
