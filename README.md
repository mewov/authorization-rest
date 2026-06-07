# Authorization Service 🔐

A user authentication and session management microservice built with **Golang (Gin)**, backed by **PostgreSQL**, and protected by a built-in **rate limiter**.

---

## ✨ Core Features & API Endpoints

Default Base URL: `http://localhost:8080`

* `GET` `/v1/status` — Checks service health.
* `POST` `/v1/auth/register` — Registers a new user (login, password, email, client, role) and issues tokens.
* `POST` `/v1/auth/login` — Authenticates user credentials and returns a JWT token pair.
* `POST` `/v1/auth/info` — Decodes the access token and returns user profile data.
* `POST` `/v1/auth/refresh` — Renews the session and issues a new access/refresh token pair.
* `POST` `/v1/auth/logout` — Terminates the active session and revokes the refresh token.

---

## ⚙️ Tech Stack

* **Language:** Go 1.25+
* **Framework:** Gin
* **Database:** PostgreSQL (stores users and active sessions)
* **Tokens:** JWT (short-lived access token ~15 min + long-lived refresh token ~7 days)
* **Middleware:** RateLimiter, structured logging (`slog`)
* **Orchestration:** Docker Compose

---

## ⚡ Quick Start

**Run Locally:**

```bash
go run cmd/server.go

```

**Run via Docker:**

```bash
docker-compose up --build

```

---

## 📄 License

**MIT License** — free to use and modify.
