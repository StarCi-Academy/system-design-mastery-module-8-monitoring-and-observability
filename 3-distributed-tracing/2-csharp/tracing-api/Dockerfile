FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS builder
WORKDIR /src
COPY tracing-api.csproj ./
RUN dotnet restore
COPY . ./
RUN dotnet publish -c Release -o /app

FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine
WORKDIR /app
ENV ASPNETCORE_URLS=http://+:3000
COPY --from=builder /app ./
EXPOSE 3000
ENTRYPOINT ["dotnet", "tracing-api.dll"]
