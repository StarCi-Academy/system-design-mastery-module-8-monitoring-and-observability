/**
 * Start the Nest Consul proxy API; bind 0.0.0.0 so it is reachable inside Docker.
 */
import {
    NestFactory,
} from "@nestjs/core"
import {
    ConfigService,
} from "@nestjs/config"
import {
    AppModule,
} from "./app.module"
import type {
    AppConfig,
} from "./config"

export async function bootstrap(): Promise<void> {
    const app = await NestFactory.create(AppModule)
    const config = app.get(ConfigService)
    const appRuntime = config.getOrThrow<AppConfig>("app")
    await app.listen(
        appRuntime.port,
        "0.0.0.0",
    )
}
