# Multi-stage build: restore + publish with the SDK, then run on the slim ASP.NET runtime.
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder
WORKDIR /src
COPY metrics-api.csproj ./
RUN dotnet restore
COPY . ./
RUN dotnet publish -c Release -o /app

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=builder /app ./
EXPOSE 3000
ENTRYPOINT ["dotnet", "metrics-api.dll"]
