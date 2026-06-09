# Multi-stage build for the ASP.NET Core consul-api proxy.
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder
WORKDIR /app
COPY consul-api.csproj ./
RUN dotnet restore
COPY . ./
RUN dotnet publish consul-api.csproj -c Release -o /out

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app
COPY --from=builder /out ./
EXPOSE 3000
ENTRYPOINT ["dotnet", "consul-api.dll"]
