# SpotSync API

Smart Parking & EV Charging Reservation platform — a centralized backend for
airports and malls to manage parking zones and the high-demand reservation of
limited EV charging spots.

**Live URL:** _add after deployment (Step 19)_

## Features

- JWT authentication with role-based access control (`driver`, `admin`)
- Parking zone management with dynamically calculated available spots
- Concurrency-safe reservations — a GORM transaction with row-level locking
  (`SELECT ... FOR UPDATE`) guarantees a zone can never be over-booked, even
  under simultaneous requests for the last spot
- Reservation lifecycle: create, list own, cancel (owner or admin), admin
  view-all
- Centralized `{success, message, data}` / `{success, message, errors}`
  response envelope; raw database errors are never exposed to clients

## Tech Stack

| Layer          | Technology                                    |
| -------------- | ---------------------------------------------- |
| Language       | Go 1.22+                                       |
| Web framework  | [Echo v4](https://echo.labstack.com/)          |
| ORM            | [GORM](https://gorm.io/) + PostgreSQL driver   |
| Validation     | go-playground/validator v10                    |
| Auth           | golang-jwt/jwt v5                              |
| Password hash  | golang.org/x/crypto/bcrypt (cost 12)           |
| Database       | PostgreSQL (NeonDB)                            |

## Architecture

Strict clean architecture — each layer only talks to the one directly below it:

```
Request
  -> router/        route -> middleware chain -> handler method
  -> handler/       bind + validate DTO, extract JWT claims, call service,
                     map errors to HTTP status, write response envelope
  -> service/       business rules (password hashing, JWT issuing, capacity
                     checks), calls repository
  -> repository/    all GORM operations: CRUD, transactions, row locks
  -> models/        GORM structs representing database tables
```

`dto/` sits alongside all layers and is the only shape ever exposed over
HTTP — GORM models never get serialized directly to JSON.

Dependency injection is wired by hand in `main.go`:
`repository.New*(db)` → `service.New*(repo)` → `handler.New*(service)` →
`router.SetupRoutes(...)`.

### Concurrency: the "EV Spot Bottleneck" fix

`repository/reservation_repository.go`'s `CreateWithLock` runs inside a
single `db.Transaction`:

1. Lock the target zone row with `clause.Locking{Strength: "UPDATE"}` —
   concurrent requests for the same zone now execute one at a time.
2. Count that zone's current `active` reservations (inside the same tx).
3. Reject with `ErrZoneFull` if `active_count >= total_capacity`, otherwise
   create the reservation.

Because the lock is held for the whole check-then-create sequence, two
simultaneous requests for the last spot can no longer both read a stale
count and both succeed.

## Project Structure

```
spotsync/
├── main.go              # DI wiring + server bootstrap
├── config/               # env loading, DB connection
├── models/                # GORM structs (User, ParkingZone, Reservation)
├── dto/                   # request/response payloads
├── repository/            # GORM data access, transactions, row locks
├── service/                # business logic
├── handler/                # HTTP handlers
├── middleware/              # JWT auth + role guard
├── router/                   # route registration
└── utils/                     # response envelope, JWT, bcrypt helpers
```

## Setup (local development)

```bash
git clone https://github.com/Arifur999/spotsync.git
cd spotsync
cp .env.example .env   # fill in real values, see below
go mod download
go run .
```

Server starts on `http://localhost:8080` (or `$PORT`). Tables are created
automatically via `AutoMigrate` on startup.

### Environment variables

| Variable          | Description                                              | Example                                                        |
| ----------------- | ---------------------------------------------------------- | --------------------------------------------------------------- |
| `PORT`            | HTTP port                                                  | `8080`                                                           |
| `DSN`             | PostgreSQL connection string                               | `postgresql://user:pass@host/db?sslmode=require`                |
| `JWT_SECRET`      | Secret used to sign JWTs                                   | a long random string                                             |
| `ALLOWED_ORIGINS` | Comma-separated CORS origins, `*` for all                 | `*`                                                               |

## API Endpoints

All routes are prefixed with `/api/v1`.

### Auth

| Method | Path            | Access | Description             |
| ------ | --------------- | ------ | ------------------------ |
| POST   | `/auth/register`| Public | Register a new user      |
| POST   | `/auth/login`   | Public | Log in, returns a JWT     |

### Parking Zones

| Method | Path         | Access      | Description                              |
| ------ | ------------ | ----------- | ------------------------------------------ |
| POST   | `/zones`     | Admin only  | Create a parking zone                      |
| GET    | `/zones`     | Public      | List all zones with dynamic `available_spots` |
| GET    | `/zones/:id` | Public      | Get a single zone with `available_spots`    |

### Reservations

| Method | Path                          | Access                | Description                              |
| ------ | ----------------------------- | ---------------------- | ------------------------------------------ |
| POST   | `/reservations`               | Authenticated           | Reserve a spot (concurrency-safe)          |
| GET    | `/reservations/my-reservations`| Authenticated          | List the caller's own reservations         |
| DELETE | `/reservations/:id`            | Authenticated (owner or admin) | Cancel a reservation                |
| GET    | `/reservations`                | Admin only              | List every reservation in the system        |

## Deployment

The repo ships a multi-stage `Dockerfile`, which Render, Railway, and Fly.io
can all build and run directly — no platform-specific build/start commands
needed.

### Render

1. New → Blueprint, point it at this repo (`render.yaml` is auto-detected),
   or New → Web Service → Environment: **Docker**.
2. Set the env vars `DSN`, `JWT_SECRET`, `ALLOWED_ORIGINS` in the dashboard
   (Render supplies `PORT` automatically).
3. Deploy — Render builds the `Dockerfile` and starts the container.

### Railway

1. New Project → Deploy from GitHub repo.
2. Railway detects the `Dockerfile` automatically and builds it.
3. Add the same env vars (`DSN`, `JWT_SECRET`, `ALLOWED_ORIGINS`) in the
   Variables tab. Railway injects `PORT` automatically.

### Fly.io

```bash
fly launch   # detects the Dockerfile, generates fly.toml
fly secrets set DSN="..." JWT_SECRET="..." ALLOWED_ORIGINS="*"
fly deploy
```

## Error Handling

Every response follows one of two shapes:

```json
{ "success": true, "message": "...", "data": { } }
```

```json
{ "success": false, "message": "...", "errors": "..." }
```

Unexpected/internal failures are logged server-side and returned to the
client as a generic `500` with no internal details attached.
