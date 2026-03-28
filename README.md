# Cinema System

 ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![CSS3](https://img.shields.io/badge/css3-%231572B6.svg?style=for-the-badge&logo=css3&logoColor=white) ![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white) ![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E) ![HTML5](https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white) 


- Project from: https://roadmap.sh/projects/movie-reservation-system

A REST API for managing a cinema workflow: users, movies, auditoriums, projections, tickets, and reservations.

The application is written in Go and uses:

- PostgreSQL for persistent storage
- Redis for movie and ticket caching
- Cookie-based sessions for authentication
- Role-based authorization for admin-only routes

## Features

- User registration, login, logout, and account deletion
- Session-based authentication using `HttpOnly` cookies
- Admin role enforcement for management endpoints
- CRUD for movies, auditoriums, projections, and ticket types
- Reservation creation with seat and ticket validation
- Pagination support for listing movies
- Redis cache for selected movie and ticket queries
- IP-based rate limiting on authentication endpoints

## Tech Stack

- Go 1.25
- net/http + `http.ServeMux` method/path patterns
- PostgreSQL 15
- Redis 7
- Docker / Docker Compose

## Project Structure

```text
cmd/cinemasys/           # Application entrypoint
internal/cache/          # Redis cache client
internal/database/       # Repositories + schema initialization
internal/domain/         # Domain entities and validations
internal/handler/        # HTTP handlers
internal/middleware/     # Auth, authorization, and rate limiting
internal/router/         # Route wiring
internal/server/         # HTTP server setup
assets/avatars/          # Allowed profile avatars
```

## How It Works

At startup, the app:

1. Loads environment variables (from `.env` if present).
2. Connects to PostgreSQL.
3. Creates required tables if they do not exist.
4. Connects to Redis and performs a health check.
5. Starts the HTTP server on the configured port.

## Data Model Overview

Main entities:

- `users`
- `sessions`
- `movies`
- `auditoriums`
- `projections`
- `tickets`
- `reservations`
- `reservation_tickets`
- `reservation_seats`

Database tables are auto-created at application startup.

## Environment Variables

### Application

- `DATABASE_URL`
	- Default: `postgres://cinemasys:secretpassword@localhost:5432/cinemasys?sslmode=disable`
- `REDIS_ADDR`
	- Default: `localhost:6379`
- `PORT`
	- Default: `8080`
- `ENV`
	- When set to `production`, session cookie is marked as `Secure=true`.

### Docker Compose Database Defaults

- `POSTGRES_USER` (default: `cinemasys`)
- `POSTGRES_PASSWORD` (default: `secretpassword`)
- `POSTGRES_DB` (default: `cinemasys`)

## Running the Project

### Option 1: Docker Compose (recommended)

```bash
docker compose up --build
```

API base URL:

```text
http://localhost:8080
```

### Option 2: Local Go Run

1. Start PostgreSQL and Redis locally.
2. Export required env vars (or create `.env`).
3. Run:

```bash
go run ./cmd/cinemasys
```

## Authentication and Authorization

### Session Cookie

- On login, the API sets a `session_id` cookie.
- Cookie is `HttpOnly` and `SameSite=Strict`.
- Session expiration is 30 minutes from login.

### Role Rules

- Regular users can use public and authenticated routes.
- Admin-only routes require both:
	- valid authenticated session
	- user role `ADMIN`

### Important Bootstrap Note

Admin routes (including user promotion) are admin-protected. In a fresh database, there is no bootstrap admin user created automatically. If you need initial admin access, you must promote a user manually in the database.

## Rate Limiting

Rate limiting is applied to:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

Behavior:

- Token-bucket style per client IP
- Initial burst capacity: 3 requests
- Refill rate: 0.05 tokens/second (approximately 3 tokens/minute)
- Exceeded requests return `429 Too Many Requests`

## API Overview

Base path: `/api/v1`

Legend:

- `Public`: no auth required
- `Auth`: requires valid `session_id` cookie
- `Admin`: requires authenticated admin user

### Auth

- `POST /auth/register` (Public)
- `POST /auth/login` (Public)
- `DELETE /auth/logout` (Auth)

### Users

- `GET /users/me` (Auth)
- `DELETE /users/me` (Auth)
- `PUT /users/{user_id}/admin` (Admin)
- `DELETE /users/{user_id}` (Admin)

### Movies

- `GET /movies` (Public)
- `GET /movies/soon` (Public)
- `GET /movies/available_now` (Public)
- `GET /movies/{movie_id}` (Public)
- `POST /movies` (Admin)
- `PUT /movies/{movie_id}` (Admin)
- `DELETE /movies/{movie_id}` (Admin)

### Projections

- `GET /movies/{movie_id}/projections` (Public)
- `POST /projections` (Admin)
- `GET /projections/{projection_id}` (Admin)
- `PUT /projections/{projection_id}` (Admin)
- `DELETE /projections/{projection_id}` (Admin)

### Auditoriums

- `GET /auditoriums` (Admin)
- `GET /auditoriums/{auditorium_id}` (Admin)
- `POST /auditoriums` (Admin)
- `PUT /auditoriums/{auditorium_id}` (Admin)
- `DELETE /auditoriums/{auditorium_id}` (Admin)

### Tickets

- `GET /tickets` (Auth)
- `POST /tickets` (Admin)
- `PUT /tickets/{ticket_id}` (Admin)
- `DELETE /tickets/{ticket_id}` (Admin)

### Reservations

- `POST /reservations` (Auth)
- `GET /reservations` (Auth)

## Request/Response Examples

All JSON requests should include:

```http
Content-Type: application/json
```

### Register

```bash
curl -i -X POST http://localhost:8080/api/v1/auth/register \
	-H "Content-Type: application/json" \
	-d '{
		"email": "user@example.com",
		"password": "StrongPass1!",
		"documentNumber": "12345678",
		"profilePicture": "/assets/avatars/batman.webp"
	}'
```

### Login

```bash
curl -i -X POST http://localhost:8080/api/v1/auth/login \
	-H "Content-Type: application/json" \
	-d '{
		"email": "user@example.com",
		"password": "StrongPass1!"
	}'
```

Store cookies in a cookie jar for subsequent requests:

```bash
curl -i -c cookies.txt -X POST http://localhost:8080/api/v1/auth/login \
	-H "Content-Type: application/json" \
	-d '{"email":"user@example.com","password":"StrongPass1!"}'
```

### Create Movie (Admin)

```bash
curl -i -b cookies.txt -X POST http://localhost:8080/api/v1/movies \
	-H "Content-Type: application/json" \
	-d '{
		"title": "Dune Part Two",
		"description": "Epic sci-fi continuation.",
		"posterImageUrl": "https://www.movieposters.com/cdn/shop/files/dune-part-two_tidfhvjl.jpg?v=1762975408&width=1680",
		"trailerUrl": "https://www.youtube.com/watch?v=Way9Dexny3w",
		"genres": ["Sci-Fi", "Action"],
		"releaseDate": "2026-03-01T20:00:00Z"
	}'
```

### List Movies with Pagination

```bash
curl -i "http://localhost:8080/api/v1/movies?limit=20&offset=0"
```

### Create Ticket Type (Admin)

```bash
curl -i -b cookies.txt -X POST http://localhost:8080/api/v1/tickets \
	-H "Content-Type: application/json" \
	-d '{
		"name": "General",
		"price": 10.50,
		"cant_seats": 1
	}'
```

### Create Reservation (Auth)

Note: the number of seats must equal the sum of `cant_seats` in provided tickets.

```bash
curl -i -b cookies.txt -X POST http://localhost:8080/api/v1/reservations \
	-H "Content-Type: application/json" \
	-d '{
		"projectionId": 1,
		"seats": [
			{"row": 2, "col": 3},
			{"row": 2, "col": 4}
		],
		"tickets": [
			{"ticketId": 1, "name": "General", "price": 10.5, "cant_seats": 2}
		]
	}'
```

## Validation Rules and Constraints

### User

- Email must have valid email format.
- Password must:
	- be at least 8 and less than 32 characters
	- include uppercase, lowercase, digit, and special character
- `documentNumber` cannot be empty.
- `profilePicture` must be one of the allowed avatar paths:
	- `/assets/avatars/batman.webp`
	- `/assets/avatars/joker.webp`
	- `/assets/avatars/spiderman.webp`
	- `/assets/avatars/dune.webp`
	- `/assets/avatars/deniro.webp`
	- `/assets/avatars/dicaprio.webp`
	- `/assets/avatars/maverick.webp`
	- `/assets/avatars/samuel.webp`
	- `/assets/avatars/travolta.webp`

### Movie

- `title` and `description` required
- `genres` required and must be valid enum values
- `posterImageUrl` and `trailerUrl` must be valid `http`/`https` URLs

### Projection

- `screeningFormat` must be `2D` or `3D`
- `language` must be `Spanish`, `Original`, or `Other`

### Ticket

- `price` must be greater than 0 on creation
- `cant_seats` must be greater than 0

### Reservation

- Total seat count must match total seats implied by included tickets
- Seat position uniqueness per projection is enforced at DB level

## Caching Behavior

Redis is used for:

- `GET /api/v1/movies/{movie_id}`
- `GET /api/v1/movies/available_now`
- `GET /api/v1/movies/soon`
- `GET /api/v1/tickets`

Cache entries are invalidated on:

- Movie update and delete operations
- Ticket create, update, and delete operations

Cache TTL is 1 hour for cached movie and ticket list/detail payloads.

## Common HTTP Status Codes

- `200 OK`
- `201 Created`
- `204 No Content`
- `400 Bad Request`
- `401 Unauthorized`
- `404 Not Found`
- `429 Too Many Requests`
- `500 Internal Server Error`