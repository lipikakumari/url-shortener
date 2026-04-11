# 🔗 URL Shortener
A production-grade URL shortener built with Go, featuring caching, message queues, and full Docker support.
## 🏗️ Architecture
```
User → Gin API → Redis Cache → Postgres
                      ↓
               RabbitMQ Queue
                      ↓
                   Worker → Postgres (click tracking)
```
## 🛠️ Tech Stack
| Technology | Purpose |
|------------|---------|
| **Go + Gin** | HTTP API server |
| **PostgreSQL** | Permanent data storage |
| **Redis** | Caching layer for fast redirects |
| **RabbitMQ** | Message queue for background processing |
| **Docker + Compose** | Containerisation |
## 📦 Project Structure
```
url-shortener/
├── main.go            # API server + handlers
├── db.go              # Postgres connection pool
├── cache.go           # Redis client
├── queries.go         # Database + cache queries
├── queue.go           # RabbitMQ publisher
├── Dockerfile         # API multistage Dockerfile
├── docker-compose.yml # All services
├── go.mod
├── go.sum
└── worker/
    ├── main.go        # Worker / consumer
    └── Dockerfile     # Worker multistage Dockerfile
```
## 🚀 Getting Started
### Prerequisites
- [Docker](https://docker.com/products/docker-desktop) installed and running
### Run everything with one command:
```bash
docker compose up --build
```
This starts all 5 services:
- Postgres on port `5432`
- Redis on port `6379`
- RabbitMQ on port `5672`
- API on port `8080`
- Worker (background)
## 📡 API Endpoints
### Shorten a URL
```
POST /shorten
```
**Request body:**
```json
{
    "url": "https://google.com"
}
```
**Response:**
```json
{
    "short_url": "http://localhost:8080/abc123"
}
```
### Redirect to original URL
```
GET /:code
```
Redirects to the original URL and tracks the click.


## ⚙️ Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://postgres:pass@localhost:5432/postgres` | Postgres connection string |
| `REDIS_URL` | `localhost:6379` | Redis address |
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection string |



## 🔄 How It Works
### Shorten URL flow:
1. `POST /shorten` receives a long URL
2. Generates a random 6-character code
3. Saves to Postgres
4. Caches in Redis immediately
5. Returns the short URL
### Redirect flow:
1. `GET /:code` receives a short code
2. Checks Redis first (cache hit → instant redirect ⚡)
3. On cache miss → queries Postgres → saves to Redis
4. Publishes click event to RabbitMQ
5. Redirects user to original URL
### Worker flow:
1. Listens to RabbitMQ `clicks` queue
2. Picks up click events
3. Increments click count in Postgres
## 🗄️ Database Schema
```sql
CREATE TABLE urls (
    code         TEXT PRIMARY KEY,
    original_url TEXT,
    clicks       INT DEFAULT 0
);
```
## 🧪 Testing
Use [Postman](https://postman.com) or curl:
```bash
# Shorten a URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
# Use the short URL
curl -L http://localhost:8080/abc123
```


