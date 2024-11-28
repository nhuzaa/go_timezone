# Toronto Time API

A Go-based API that provides the current time in Toronto and logs all requests to a MySQL database.

## Features

- Get current Toronto time
- Log all time requests to MySQL database
- Retrieve time request history
- Docker containerization
- Structured logging

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. Clone the repository
2. Run the application:

   ```bash
   docker-compose up --build
   ```

## API Endpoints

- `GET /current-time`: Returns the current time in Toronto
- `GET /time-logs`: Returns the history of time requests

## Environment Variables

The following environment variables can be configured in docker-compose.yml:

- `DB_HOST`: MySQL host
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

## Development

To run tests:
