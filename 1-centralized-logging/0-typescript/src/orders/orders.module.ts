/**
 * Orders module — registers the controller and service for log/error demos.
 */
import { Module } from "@nestjs/common"
import { OrdersController } from "./orders.controller"
import { OrdersService } from "./orders.service"

@Module({
    controllers: [OrdersController],
    providers: [OrdersService],
})
export class OrdersModule {}
