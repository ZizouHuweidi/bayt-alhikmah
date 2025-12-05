.PHONY: run docker-up db-shell migrate-add migrate-up test

run:
	dotnet run --project src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj

docker-up:
	docker-compose up -d

db-shell:
	docker exec -it bayt-alhikmah-postgres-1 psql -U postgres -d bayt_alhikmah

migrate-add:
	@read -p "Enter migration name: " name; \
	dotnet ef migrations add $$name --project src/BaytAlHikmah.Infrastructure/BaytAlHikmah.Infrastructure.csproj --startup-project src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj

migrate-up:
	dotnet ef database update --project src/BaytAlHikmah.Infrastructure/BaytAlHikmah.Infrastructure.csproj --startup-project src/BaytAlHikmah.Api/BaytAlHikmah.Api.csproj

test:
	dotnet test
