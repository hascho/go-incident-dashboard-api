# Distributed Incident Management & Notification System

A production-grade Incident Management API built with **Go**, featuring a high-performance background worker system. This project demonstrates the transition from a standard CRUD API to a distributed architecture using **PostgreSQL** for persistence and **Redis** for real-time event signaling.



## 🚀 System Architecture

The system employs a **Hybrid Producer-Consumer** pattern designed for high availability and low latency:

* **The API (Producer):** Handles incident creation and performs a "dual-write" strategy. It persists the incident and a corresponding "Notification Job" to PostgreSQL, then immediately publishes the Job ID to Redis.
* **The Worker (Consumer):** Acts as a real-time listener. It subscribes to a Redis channel for instant triggers but also runs a "Safety Poll" every 30 seconds to catch any jobs missed due to network fluctuations.
* **Decoupled Logic:** The worker handles the heavy lifting (simulated notifications), ensuring the API remains fast and responsive for the end-user.

## 🛠 Tech Stack

* **Language:** Go (Golang)
* **Web Framework:** Gin Gonic
* **Database:** PostgreSQL (Relational storage & Job Queueing)
* **Message Broker:** Redis (Pub/Sub for real-time triggers)
* **Migrations:** Goose
* **Logging:** Zerolog (Structured JSON logging)
* **Containerization:** Docker & Docker Compose

## 🧠 Key Learnings & Engineering Patterns

### 1. The Interface-Driven Pattern
To keep the code maintainable, I implemented the **TaskQueue** interface. This allows the application to switch from Redis to other brokers (like NATS or RabbitMQ) without changing a single line of business logic in the service layer.

### 2. Defensive Engineering (The Safety Net)
I implemented a hybrid worker that combines **Pub/Sub (Speed)** with **Polling (Reliability)**. This ensures 100% task completion even if the message broker restarts or fails briefly.

### 3. Database Schema Evolution
The project uses **Goose** for versioned migrations. This allowed for seamless updates to the schema, such as adding the `notification_status` column to existing tables while preserving data integrity.

---

## 💻 Common Commands

### Infrastructure
```bash
# Start Postgres and Redis
docker-compose up -d

# Check running containers
docker ps
```

### Database Migrations (Goose)
```bash
# Create a new migration file
goose -dir internal/db/migrations create <name_of_migration> sql

# Apply migrations
GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=user password=password dbname=incidentdb sslmode=disable" goose -dir internal/db/migrations up
```

### Running Services
```bash
# Run the API Server
go run cmd/server/main.go

# Run the Background Worker
go run cmd/worker/main.go
```

### Debugging Tools
```bash
# Monitor real-time Redis signals
docker exec -it $(docker ps -qf "name=redis") redis-cli
SUBSCRIBE notification_jobs

# Verify results in PostgreSQL
psql -U user -d incidentdb
SELECT title, status, notification_status FROM incidents ORDER BY created_at DESC LIMIT 5;
```

---

## 📈 Challenges Overcome

* **JSON Case Sensitivity:** Encountered an issue where PostgreSQL JSONB keys were capitalized by Go's default Marshaller. Resolved this by implementing explicit struct tags (`json:"title"`) to ensure consistency between the Go application and the database.
* **Graceful Shutdowns:** Implemented OS signal handling (SIGINT/SIGTERM) in both the API and Worker to ensure that database connections are closed cleanly and current jobs are finished before the process exits.
* **Dependency Injection:** Managed complex dependencies (multiple repositories and queues) by injecting them through constructors, making the code more testable and modular.

---

## 🔄 Project Evolution
1.  **Stage 1 & 2:** Standard REST API with CRUD operations.
2.  **Stage 3:** Added a persistent Job Queue table in Postgres to handle retries.
3.  **Stage 4:** Integrated Redis Pub/Sub to move from 5-second polling to millisecond-latency triggers.
4.  **Stage 5:** Implemented a full-circle status update, where workers report success back to the primary incident table.

## 📖 API Reference

All requests should be sent to `http://localhost:8080`.

### 1. Create an Incident
This endpoint triggers the background worker via Redis.

**Request:**
```bash
curl -X POST http://localhost:8080/incidents \
-H "Content-Type: application/json" \
-d '{
  "title": "Service Timeout",
  "description": "API Gateway is timing out",
  "status": "open",
  "severity": "high",
  "team": "DevOps"
}'
```

**Response (201 Created):**
```json
{
  "id": "878b6f82-5075-4b03-828f-7e1fe189a5e8",
  "title": "Service Timeout",
  "status": "open",
  "severity": "high",
  "team": "DevOps"
}
```

---

### 2. List All Incidents
Returns a list of all incidents including their background notification status.

**Request:**
```bash
curl -X GET http://localhost:8080/incidents
```

---

### 3. Get Incident by ID
Retrieve details for a specific incident.

**Request:**
```bash
curl -X GET http://localhost:8080/incidents/878b6f82-5075-4b03-828f-7e1fe189a5e8
```

---

### 4. Update an Incident
Update the status or description of an existing incident.

**Request:**
```bash
curl -X PATCH http://localhost:8080/incidents/878b6f82-5075-4b03-828f-7e1fe189a5e8 \
-H "Content-Type: application/json" \
-d '{
  "status": "resolved",
  "description": "Issue fixed by scaling the gateway"
}'
```

---

### 5. Delete an Incident
Remove an incident from the database.

**Request:**
```bash
curl -X DELETE http://localhost:8080/incidents/878b6f82-5075-4b03-828f-7e1fe189a5e8
```

---

## 🛠 Notification Statuses

The `notification_status` field in the database tracks the background worker's progress:

| Status | Description |
| :--- | :--- |
| `pending` | The incident was created, and a job is queued. |
| `sent` | The background worker has successfully processed the notification. |
| `failed` | The worker exhausted retry attempts (Check worker logs). |