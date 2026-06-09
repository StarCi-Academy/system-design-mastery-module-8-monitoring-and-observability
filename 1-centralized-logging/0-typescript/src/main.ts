/**
 * Node entry — import the Sentry instrument first, then invoke bootstrap.
 */
// Sentry must be imported first to hook into the runtime before any other module.
import "./instrument"
import { bootstrap } from "./bootstrap"

void bootstrap()
