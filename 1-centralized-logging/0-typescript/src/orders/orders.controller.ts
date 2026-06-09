/**
 * Orders controller — routes delegate to the service.
 */
import { Controller, Get } from "@nestjs/common"
import { OrdersService } from "./orders.service"

@Controller("api/orders")
export class OrdersController {
    constructor(private readonly ordersService: OrdersService) {}

    /**
     * GET /api/orders — accepted order; logs INFO to Loki for the query demo.
     */
    @Get()
    accept(): { status: string; trace: string } {
        return this.ordersService.accept()
    }

    /**
     * GET /api/orders/error — deliberate failure; Sentry captures the exception,
     * Winston logs ERROR to Loki.
     */
    @Get("error")
    error(): never {
        return this.ordersService.fail()
    }
}
