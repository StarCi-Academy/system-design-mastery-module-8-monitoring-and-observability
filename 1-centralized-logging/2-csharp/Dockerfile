# Multi-stage build for the ASP.NET Core logging-api.
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder
WORKDIR /src
COPY logging-api.csproj ./
RUN dotnet restore
COPY Program.cs ./
RUN dotnet publish -c Release -o /app

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=builder /app ./
EXPOSE 3000
ENTRYPOINT ["dotnet", "logging-api.dll"]
