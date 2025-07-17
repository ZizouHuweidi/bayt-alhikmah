FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build

WORKDIR /app

COPY src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj src/BaytAlHikmah.Api/
RUN dotnet restore src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj

COPY . .

RUN dotnet publish src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj -c Release -o out

FROM mcr.microsoft.com/dotnet/aspnet:9.0

WORKDIR /app

COPY --from=build /app/out .

EXPOSE 8080

ENV ASPNETCORE_URLS=http://+:8080

ENTRYPOINT ["dotnet", "BaytAlHikmah.Api.dll"]