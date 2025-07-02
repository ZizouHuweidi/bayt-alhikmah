# Bayt al-Hikmah

Bayt al-Hikmah is a modern platform for managing and engaging with knowledge sources across formats—including books, academic papers, podcasts, videos, and articles. Its core purpose is to help users collect, organize, annotate, and track their engagement with these sources, while optionally contributing to a shared public knowledge space.

## Features

- **Structured knowledge source tracking:** Track metadata like title, author, tags, type, and progress.
- **Notes and annotations:** Create private and public notes linked to specific sources, pages, or timestamps.
- **Reading and knowledge timeline:** Visualize your engagement over time.
- **Lists and collections:** Organize sources by theme or intent.
- **Tags, topics, and taxonomies:** Contextualize your knowledge with flexible organization.
- **Ratings and reviews:** Evaluate and reflect on sources.
- **User profiles:** Showcase your reading history, contributions, and knowledge areas.

## Getting Started

### Prerequisites

- [.NET 8 SDK](https://dotnet.microsoft.com/download/dotnet/8.0)
- [Docker](https://www.docker.com/get-started)

### Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/zizouhuweidi/bayt-alhikmah.git
   cd bayt-alhikmah
   ```

2. **Build and run the application using Docker Compose:**

   ```sh
   docker-compose up --build
   ```

   The API will be available at `http://localhost:8080`.

## Usage

- **API Documentation:** Once the application is running, you can access the Swagger UI for API documentation at `http://localhost:8080/swagger`.
- **Admin User:** The first user to register will be automatically assigned the `Admin` role.
- **Admin Endpoints:** Admin-only endpoints are available at `/admin` for managing users and other administrative tasks.

## Development

### Running Locally

To run the application locally without Docker:

1. **Start the PostgreSQL database:**

   ```sh
   docker-compose up -d postgres
   ```

2. **Run the API:**
   ```sh
   dotnet run --project src/BaytAlHikmah.Api
   ```
