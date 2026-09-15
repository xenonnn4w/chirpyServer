# Chirpy

A small Twitter-style REST API built with Go and PostgreSQL. It supports user accounts, authenticated chirp creation and deletion, token refresh, and a simple web interface.

## Architecture

![Detailed Chirpy system architecture](docs/architecture/chirpy-system-architecture.drawio.svg)

## Stack

- Go 1.24 with `net/http`
- PostgreSQL with sqlc-generated queries and goose migrations
- Argon2id password hashing
- HS256 JWT access tokens and opaque refresh tokens

## Run locally

Requirements: Go 1.24+, PostgreSQL, and [goose](https://github.com/pressly/goose).

Create a `.env` file:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=replace-with-a-random-secret
POLKA_KEY=replace-with-a-webhook-key
PLATFORM=dev
```

Apply the database migrations and start the server:

```bash
set -a
source .env
set +a
goose -dir sql/schema postgres "$DB_URL" up
go run .
```

The API runs at `http://localhost:8080`; the web interface is available at `http://localhost:8080/app/`.

## Core API

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `POST` | `/api/users` | Create an account |
| `POST` | `/api/login` | Receive access and refresh tokens |
| `POST` | `/api/refresh` | Issue a new access token |
| `POST` | `/api/revoke` | Revoke a refresh token |
| `GET` / `POST` | `/api/chirps` | List or create chirps |
| `GET` / `DELETE` | `/api/chirps/{chirpID}` | Read a chirp or delete an owned chirp |
| `PUT` | `/api/users` | Update the authenticated user |
| `POST` | `/api/polka/webhooks` | Process Chirpy Red upgrades |

Protected endpoints use `Authorization: Bearer <token>`. The webhook uses `Authorization: ApiKey <key>`.

## Security review

The separate [Chirpy Security Lab](https://github.com/xenonnn4w/chirpy-security-lab) documents a reproducible JWT validation review, hardened policy, fuzz tests, and object-authorization regression tests.
