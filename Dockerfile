# 1. Build Stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /app

# Copy solution and project files
COPY *.sln .
COPY src/*/*.csproj ./src/
COPY tests/*/*.csproj ./tests/

# Restore dependencies for all projects
RUN dotnet restore

# Copy the rest of the source code
COPY . .

# Publish the API project
WORKDIR /app/src/BaytAlHikmah.Api
RUN dotnet publish -c Release -o /app/out

# 2. Final Stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app/out .

EXPOSE 8080
ENV ASPNETCORE_URLS=http://+:8080

ENTRYPOINT ["dotnet", "BaytAlHikmah.Api.dll"]
