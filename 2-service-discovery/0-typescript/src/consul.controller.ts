/**
 * Consul proxy controller — register / health / deregister via the Consul Agent HTTP API.
 */
import {
    Body,
    Controller,
    Get,
    Param,
    Post,
} from "@nestjs/common"
import {
    ConfigService,
} from "@nestjs/config"
import axios from "axios"
import type {
    ConsulConfig,
} from "./config"

@Controller("consul")
export class ConsulController {
    private readonly baseUrl: string

    public constructor(private readonly configService: ConfigService) {
        const consul = this.configService.getOrThrow<ConsulConfig>("consul")
        this.baseUrl = consul.baseUrl
    }

    /**
     * Register a service instance into the Consul catalog via the Agent HTTP API.
     */
    @Post("register")
    public async register(@Body() body: { id: string; name: string; address: string; port: number }): Promise<{ status: string }> {
        // Consul requires PascalCase fields: ID, Name, Address, Port.
        await axios.put(`${this.baseUrl}/v1/agent/service/register`, {
            ID: body.id,
            Name: body.name,
            Address: body.address,
            Port: body.port,
        })
        return {
            status: "registered",
        }
    }

    /**
     * List passing instances of a service by logical name.
     */
    @Get("health/:service")
    public async health(@Param("service") service: string): Promise<unknown> {
        const response = await axios.get(`${this.baseUrl}/v1/health/service/${service}`)
        return response.data
    }

    /**
     * Deregister a service instance from the Agent by its exact instance ID.
     */
    @Post("deregister/:id")
    public async deregister(@Param("id") id: string): Promise<{ status: string }> {
        await axios.put(`${this.baseUrl}/v1/agent/service/deregister/${id}`)
        return {
            status: "deregistered",
        }
    }
}
