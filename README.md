# 💬 BukChat - Real-Time Chat & Collaboration Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.2-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

**BukChat** is a high-performance, real-time messaging application engineered with a **Go Clean Architecture backend**, a **React 18 TypeScript frontend**, and an event-driven **Gorilla WebSocket Hub** engine with structured `slog` logging and Pub/Sub abstraction.

---

## 📌 Table of Contents

- [System Architecture](#-system-architecture)
- [Tech Stack](#-tech-stack)
- [Backend Clean Architecture](#-backend-clean-architecture)
- [Database Schema & ER Diagram](#-database-schema--er-diagram)
- [Real-Time WebSocket Engine](#-real-time-websocket-engine)
- [Authentication & Security](#-authentication--security)
- [API Endpoint Reference](#-api-endpoint-reference)
- [DevOps & Deployment Pipeline](#-devops--deployment-pipeline)
- [Local Development Setup](#-local-development-setup)

---

## 🏗 System Architecture

The application follows an enterprise web architecture with client-side rendering, reverse-proxy routing, stateless RESTful APIs, and stateful WebSocket channels:

```
+-------------------------------------------------------------------------+
|                              CLIENT LAYER                               |
|            Browser (React 18 SPA - Vite + TypeScript + Tailwind)        |
+-------------------------------------------------------------------------+
                                   |
                       HTTP / WSS  |  (SSL Termination & Path Routing)
                                   v
+-------------------------------------------------------------------------+
|                          NGINX REVERSE PROXY                            |
+-------------------------------------------------------------------------+
                    /                                 \
      Proxy /v1    /                                   \   Proxy /ws
                  v                                     v
+-----------------------------------+   +---------------------------------+
|          GIN HTTP ROUTER          |   |      WEBSOCKET HUB GOROUTINE    |
|   (/v1/auth, /v1/users, /v1/chat) |   |          (/ws/:roomId)          |
+-----------------------------------+   +---------------------------------+
                  \                                     /
                   \--->    BACKEND CLEAN ARCHITECTURE <---/
                         (Controllers -> Usecases -> Repositories)
                                       |
                                       | SQLx Queries
                                       v
+-------------------------------------------------------------------------+
|                            PERSISTENCE LAYER                            |
|             PostgreSQL 14 (Users, Rooms, Messages, Friends)             |
+-------------------------------------------------------------------------+
```

---

## 🛠 Tech Stack

### **Frontend**
- **Core Framework:** [React 18](https://react.dev/) + [TypeScript](https://www.typescriptlang.org/)
- **Build Tool:** [Vite](https://vitejs.dev/)
- **Styling & UI:** [Tailwind CSS](https://tailwindcss.com/), Radix UI, Framer Motion, Lucide Icons
- **Real-Time Client:** [`react-use-websocket`](https://www.npmjs.com/package/react-use-websocket)
- **Production Web Server:** [Nginx](https://www.nginx.com/) (Multi-stage Docker container)

### **Backend**
- **Language:** [Go (Golang 1.21+)](https://go.dev/) (`github.com/bukharney/bukchat`)
- **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
- **WebSocket Protocol:** [Gorilla WebSocket](https://github.com/gorilla/websocket)
- **Database Driver / ORM Layer:** [`sqlx`](https://github.com/jmoiron/sqlx) + [`lib/pq`](https://github.com/lib/pq)
- **Authentication & Security:** JWT (`golang-jwt/jwt/v4`), Password Hashing (`golang.org/x/crypto/bcrypt`)
- **Observability & Logging:** Go stdlib `slog` (Structured JSON/Text logging)

### **Database & Infrastructure**
- **Database:** PostgreSQL 14 Alpine
- **Containerization:** Docker & Docker Compose
- **Host Reverse Proxy:** Nginx Reverse Proxy
- **CI/CD Pipeline:** GitHub Actions ([`.github/workflows/deploy-homelab.yml`](./.github/workflows/deploy-homelab.yml)) publishing multi-arch (ARM64/x86_64) images to GitHub Container Registry (GHCR)

---

## 🏛 Backend Clean Architecture

The backend repository ([`backend`](./backend)) implements Domain-Driven Clean Architecture with full `context.Context` propagation to ensure separation of concerns, testability, and maintainability:

```
backend/
├── configs/            # Application environment configurations
├── database/           # PostgreSQL connection & schema initialization
├── server/
│   ├── server.go       # Server initialization, CORS & Graceful Shutdown
│   ├── handler.go      # Route group mapping (/v1, /ws)
│   └── ws/             # Gorilla WebSocket Hub, Client Goroutines & PubSub Adapter
├── modules/
│   ├── controllers/    # HTTP Request Binding & Response Formatting
│   ├── usecases/       # Core Business Logic & Orchestration
│   ├── repositories/   # SQL Database Execution (SQLx Context Queries)
│   └── entities/       # Domain Models, Interfaces, Data Contracts
├── middlewares/        # JWT Authentication & Authorization Middlewares
├── pkg/
│   ├── apperrors/      # Standardized Application Error Definitions
│   └── logger/         # Structured Slog Logger Initialization
└── utils/              # Helper utilities (JSON responses, connection builders)
```

1. **Controllers ([`backend/modules/controllers`](./backend/modules/controllers)):** Parse HTTP payloads, validate inputs, extract JWT claims, pass `c.Request.Context()`, and format standardized API responses.
2. **Usecases ([`backend/modules/usecases`](./backend/modules/usecases)):** Enforce business rules (e.g. friend request acceptance creating chat rooms, user status changes broadcasting live notifications) with mockable unit tests.
3. **Repositories ([`backend/modules/repositories`](./backend/modules/repositories)):** Execute context-aware SQL queries against PostgreSQL using `sqlx`.
4. **Entities ([`backend/modules/entities`](./backend/modules/entities)):** Interface contracts binding layers together through Go interfaces.

---

## 🗄 Database Schema & ER Diagram

BukChat's relational schema handles user identity, room management, message persistence, and social friendships:

```
+------------------+       +------------------+       +------------------+
|      USERS       |       |   USERS_ROOMS    |       |      ROOMS       |
+------------------+       +------------------+       +------------------+
| id (PK)   SERIAL |1    N | user_id (PK, FK) | N    1| id (PK)   SERIAL |<---+
| username  VARCHAR|<----->| room_id (PK, FK) |<----->| name     VARCHAR |    |
| email     VARCHAR|       +------------------+       | created_at  STAMP|    | 1
| password  VARCHAR|                                  +------------------+    |
| created_at  STAMP|                                    | 1                   |
+------------------+                                    |                     | N
  | 1         | 1                                       v N                   | (FK)
  |           +-----------------------+               +------------------+    |
  v N                                 v N             |     MESSAGES     |    |
+------------------+                                  +------------------+    |
|     FRIENDS      |                                  | id (PK)   SERIAL |    |
+------------------+                                  | room_id (FK) INT |    |
| from_user_id (PK)|                                  | user_id (FK) INT |<---+--+
| to_user_id (PK)  |                                  | message     TEXT |    |  | 1 (author)
| status       INT |                                  | created_at  STAMP|    |  |
| room_id (FK) INT |----------------(FK to ROOMS)-----------------------------+  |
+------------------+                                                             |
  |                                                                              |
  +------------------------------------ N ---------------------------------------+
```

### Table Definitions
- **`users`**: User account credentials, email, and password hash (Bcrypt).
- **`rooms`**: Chat channels / direct message instances.
- **`messages`**: Historical log of messages belonging to a room and author.
- **`users_rooms`**: Junction table mapping many-to-many relationships between users and rooms.
- **`friends`**: Represents social connections with friend request status states (`0: pending`, `1: accepted`, `2: rejected`) and a linked direct chat `room_id`.

---

## ⚡ Real-Time WebSocket Engine

BukChat's real-time messaging engine ([`backend/server/ws`](./backend/server/ws)) uses Go's low-overhead concurrency primitives (**Goroutines, Channels, and RWMutexes**) to support multi-room broadcasting, direct user notifications, and live presence tracking without external dependencies.

---

### 1. Connection Upgrade & Authentication Flow

Standard browser `WebSocket` APIs cannot attach custom HTTP headers during the initial handshake. BukChat authenticates connections cleanly via a query parameter during protocol upgrade:

```
+----------------+        1. HTTP GET /ws/:roomId?token=<jwt>        +------------------+
| CLIENT BROWSER | -------------------------------------------------> | GIN WS ROUTER    |
+----------------+                                                   +------------------+
        |                                                                      |
        |                                                            2. Verify JWT Claims
        |                                                            3. Upgrade Connection
        v                                                                      v
+----------------+                      4. Spawn                      +------------------+
| WEBSOCKET CONN | <------------------------------------------------- | CLIENT INSTANCE  |
+----------------+         go client.Read() & go client.Write()       +------------------+
```

1. **Endpoint:** `GET /ws/:roomId?token=<JWT>` handled by [`ServeWS`](./backend/server/ws/client.go).
2. **Authentication:** Validates the JWT string using `middlewares.GetUserToken(tk)` before establishing the socket connection.
3. **Connection Upgrade:** Upgrades the standard HTTP GET request to full-duplex WebSocket using Gorilla's `upgrader.Upgrade()`.
4. **Client Binding:** Instantiates a `Client` struct containing the connection, room ID, user identity claims, and a thread-safe message buffer (`chan Message, 256`).

---

### 2. The `Hub` & `Client` Concurrency Model

A single central `Hub` goroutine manages client lifecycle state and message routing across rooms:

```
                      +--------------------------------------------------+
                      |                  WEBSOCKET HUB                   |
                      |                                                  |
 Client Connect       |  - clients:     map[roomId]map[*Client]bool      |  Client Disconnect
--------------------->|  - userClients: map[userID]map[*Client]bool      |-------------------->
  (register chan)     |  - mu:          sync.RWMutex                     |  (unregister chan)
                      +--------------------------------------------------+
                                        |                |
                       Broadcast (room) |                | Notification (user)
                                        v                v
                      +-------------------+            +-------------------+
                      | Room Broadcast    |            | User Send         |
                      | Loop (clients)    |            | Loop (userClients)|
                      +-------------------+            +-------------------+
                                        |                |
                                        +--------+-------+
                                                 | Pushes to client.send chan
                                                 v
                                    +--------------------------+
                                    | Client.Write() Goroutine |
                                    |  (Writes JSON to WS Conn)|
                                    +--------------------------+
```

#### State Tracking & Thread Safety:
- **`clients` Map (`roomId -> map[*Client]bool`):** Tracks active socket connections per chat room for room-wide broadcasting.
- **`userClients` Map (`userID -> map[*Client]bool`):** Tracks connections per user across multiple tabs/devices for targeted notifications and presence calculation.
- **`sync.RWMutex`:** Protects map mutations during concurrent connection joins, leaves, and message routing.
- **`PubSubAdapter` Interface ([`pubsub.go`](./backend/server/ws/pubsub.go)):** Provides a clean abstraction layer (`InMemoryPubSub` default) allowing seamless transition to Redis Pub/Sub for horizontal multi-instance scaling.

---

### 3. Dual-Goroutine Client Execution

Each connected WebSocket client spawns two dedicated lightweight Goroutines:

* **`Read()` Goroutine:**
  - Continuously listens for incoming JSON payloads from the client connection.
  - Attaches sender credentials (`c.User.Id`) to prevent client-side identity spoofing.
  - Configures connection read limit (`512 KB`) and Pong deadlines (`60s`).
  - **Message Processing:** If the event is standard text (`chat`), it persists the payload to PostgreSQL via `ChatRepository.SendMessage(ctx)` before dispatching it to `hub.broadcast`.
  - **Ephemeral Events:** Ephemeral events like `typing` bypass database persistence and are immediately forwarded to the room broadcast channel.
* **`Write()` Goroutine:**
  - Spawns a Ping ticker (`54s` interval) to send keep-alive frames (`websocket.PingMessage`).
  - Consumes outgoing `Message` objects from the client's buffered `send` channel (`chan Message, 256`).
  - Writes JSON frames down the active WebSocket connection with write deadlines (`10s`).

---

### 4. Event Protocol & Payload Specifications

| Event Type | Persisted to DB? | Target Audience | Description |
| :--- | :---: | :--- | :--- |
| **`chat`** | ✅ Yes | Room Subscribers (`clients[roomId]`) | Standard text message sent in a room. Saved to `messages` table before fan-out. |
| **`typing`** | ❌ Ephemeral | Room Subscribers (`clients[roomId]`) | Live user typing indicator. Dispatched in real time without DB overhead. |
| **`user_status`** | ❌ Live State | All Clients (`notifications`) | Dispatched when a user joins their first connection (`online: true`) or closes their last socket (`online: false`). |
| **`new_message`** | ❌ Notification | Target User (`userClients[userID]`) | Cross-room toast notification alerting offline/other-room users of incoming messages. |
| **`system`** | ❌ Ephemeral | Room Subscribers (`clients[roomId]`) | Server-generated events (e.g. `"User X joined the room"`). |

---

### 5. Non-Blocking Fan-Out & Slow Client Protection

To prevent a lagging client or dropped network connection from blocking the central `Hub` loop and freezing room broadcasts for all users, BukChat uses **non-blocking channel writes with `select` fallbacks**:

```go
select {
case client.send <- message:
    // Message successfully queued in client's buffer
default:
    // Client buffer is full (256 messages); queue stale client for thread-safe removal
    stale = append(stale, client)
}
```

Stale clients are subsequently cleaned up under a write lock (`h.mu.Lock()`), closing un-responsive connections safely without blocking live broadcasts.

---

## 🔐 Authentication & Security

1. **REST API Authorization:**
   - Stateless JWT tokens signed with HMAC SHA-256 (`HS256`).
   - Authorization Header: `Bearer <token>` extracted via `middlewares.JwtAuthentication()`.
2. **WebSocket Authentication:**
   - WebSockets handshakes pass JWT in query parameter: `GET /ws/:roomId?token=<jwt_token>`.
   - Middleware extracts claims and binds user ID & username to the `Client` struct before upgrading HTTP connection to WebSocket.
3. **Password Security:**
   - Passwords hashed using `bcrypt` before database write.
4. **SQL Injection Prevention:**
   - Prepared statements and parameterized queries enforced by `sqlx`.

---

## 📋 API Endpoint Reference

### Auth Endpoints (`/v1/auth`)
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/v1/auth/login` | Authenticate user credentials & return JWT | ❌ No |
| `GET` | `/v1/auth/auth-test` | Test token validity | ✅ Yes (JWT) |
| `GET` | `/v1/auth/refresh-token` | Obtain new access token | ✅ Yes (JWT) |

### User & Friend Endpoints (`/v1/users`)
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/v1/users/` | Register new user account | ❌ No |
| `GET` | `/v1/users/` | Get authenticated user profile | ✅ Yes (JWT) |
| `PATCH` | `/v1/users/` | Change password | ✅ Yes (JWT) |
| `DELETE` | `/v1/users/` | Delete user account | ✅ Yes (JWT) |
| `GET` | `/v1/users/friends` | List accepted friends | ✅ Yes (JWT) |
| `GET` | `/v1/users/friends-request` | List pending friend requests | ✅ Yes (JWT) |
| `POST` | `/v1/users/add-friend` | Send friend request | ✅ Yes (JWT) |
| `POST` | `/v1/users/reject-friend` | Reject/cancel friend request | ✅ Yes (JWT) |

### Chat & WebSocket Endpoints (`/v1/chat` & `/ws`)
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/v1/chat/` | Create new chat room | ✅ Yes (JWT) |
| `GET` | `/v1/chat/:roomId` | Fetch message history for a room | ✅ Yes (JWT) |
| `GET` | `/ws/:roomId` | Upgrade HTTP to WebSocket connection | ✅ Query (`?token=`) |

---

## 🚢 DevOps & Deployment Pipeline

BukChat features automated CI/CD using **GitHub Actions** ([`.github/workflows/deploy-homelab.yml`](./.github/workflows/deploy-homelab.yml)) building multi-architecture Docker images and publishing them to **GitHub Container Registry (GHCR)** for deployment to a local Home Lab server:

```
+-----------+            +-------------------------+            +------------------------------------+            +---------------------------------+
| DEVELOPER |            | GITHUB ACTIONS WORKFLOW |            | GITHUB CONTAINER REGISTRY (GHCR)   |            | HOME LAB SERVER (ARM64/x86_64) |
+-----------+            +-------------------------+            +------------------------------------+            +---------------------------------+
      |                               |                                           |                                           |
      |--- 1. Git push main --------->|                                           |                                           |
      |                               |--- 2. Setup QEMU & Buildx (ARM64) ------->|                                           |
      |                               |--- 3. Push client & server images ------->|                                           |
      |                                                                           |<--- 4. Pull latest images (Docker Compose)|
      |                                                                           |                                           |--- 5. Run docker compose up -d
```

- **Multi-Architecture Build:** Uses QEMU and Docker Buildx to cross-compile ARM64 & x86_64 container images suitable for Home Lab single-board computers (Raspberry Pi, Turing Pi, Apple Silicon nodes) or x86 servers.
- **Container Registry:** Images are securely hosted on `ghcr.io/bukharney/bukchat/client:latest` and `ghcr.io/bukharney/bukchat/server:latest`.
- **Reverse Proxy:** Host Nginx reverse proxy handles SSL termination and proxies requests to local container ports (`8080` backend, `4173`/`5173` frontend).

---

## 💻 Local Development Setup

### Prerequisites
- [Go 1.21+](https://go.dev/doc/install)
- [Node.js 18+](https://nodejs.org/) & `pnpm` / `npm`
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- PostgreSQL 14 (if running without Docker)

---

### Option 1: Quickstart with Docker Compose (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/bukharney/BukChat.git
   cd BukChat
   ```

2. **Create environment file `.env`:**
   ```env
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_DB=bukchat
   JWT_SECRET=supersecretkey
   LOG_LEVEL=info
   LOG_FORMAT=text
   ```

3. **Spin up entire stack (Database, Backend, Frontend):**
   ```bash
   docker-compose up --build -d
   ```
   - **Frontend:** `http://localhost:4173`
   - **Backend API:** `http://localhost:8080`
   - **PostgreSQL:** `localhost:5432`

---

### Option 2: Running Standalone (Local Development)

#### 1. Start PostgreSQL Database
```bash
docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=bukchat -p 5432:5432 -d postgres:14-alpine
```

#### 2. Start Backend Server
```bash
cd backend
go run main.go
```
*The database tables will be automatically created on startup by `database/db.go`.*

#### 3. Start Frontend Client
```bash
cd frontend
pnpm install
pnpm dev
```
*Frontend runs at `http://localhost:5173`.*

#### 4. Run Backend Unit Tests
```bash
cd backend
go test -v ./...
```
> **Note on Race Detector:** Running `go test -race ./...` on Windows requires CGO and a C compiler (`gcc`) added to `%PATH%`, or can be run inside Docker/WSL (`docker run --rm -v "${PWD}:/app" -w /app golang:1.21 go test -race ./...`).

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---
*Authored & Maintained by [Bukharney](https://github.com/bukharney)* 🚀