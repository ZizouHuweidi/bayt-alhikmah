FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /app

# Copy solution and project files
COPY BaytAlHikmah.sln ./
COPY src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj src/BaytAlHikmah.Api/
COPY src/BaytAlHikmah.Application/BaytAlHikmah.Application.csproj src/BaytAlHikmah.Application/
COPY src/BaytAlHikmah.Domain/BaytAlHikmah.Domain.csproj src/BaytAlHikmah.Domain/
COPY src/BaytAlHikmah.Infrastructure/BaytAlHikmah.Infrastructure.csproj src/BaytAlHikmah.Infrastructure/
COPY global.json ./
COPY Directory.Packages.props ./

# Restore dependencies
RUN dotnet restore

# Copy the rest of the code
COPY src/ src/

# Build and publish
WORKDIR /app/src/BaytAlHikmah.Api
RUN dotnet publish -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:10.0
WORKDIR /app
COPY --from=build /app/publish .
EXPOSE 8080
ENTRYPOINT ["dotnet", "BaytAlHikmah.Api.dll"]
