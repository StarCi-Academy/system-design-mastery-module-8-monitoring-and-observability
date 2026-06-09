/**
 * App runtime config — `app` namespace; reads process.env in the factory only.
 */
import {
    registerAs,
} from "@nestjs/config"

export interface AppConfig {
    port: number
}

export const appConfig = registerAs(
    "app",
    (): AppConfig => ({
        port: Number(process.env.PORT) || 3000,
    }),
)
