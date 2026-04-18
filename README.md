🔗 URL Shortener
A production-grade URL shortener built with Go, featuring authentication, rate limiting, caching, message queues, and full Docker support.
🏗️ Architecture
User → Rate Limiter → JWT Auth → Gin API → Redis Cache → Postgres
                                     ↓
                               RabbitMQ Queue
                                     ↓
                              Consumer → Postgres (click tracking)

🛠️ Tech Stack
Technology          Purpose
Go + Gin.           HTTP API server
PostgreSQL.         Permanent data storage
Redis               Caching layer + rate limiting
RabbitMQ            Message queue for background processing
JWT                 Authentication tokens
bcrypt              Password hashing
Docker + Compose    Containerisation

📦 Project Structure
url-shortener/
├── main.go            # API server + handlers
├── db.go              # Postgres connection pool
├── cache.go           # Redis client
├── queries.go         # Database + cache queries
├── queue.go           # RabbitMQ publisher
├── middleware.go      # Rate limiting + JWT middleware
├── auth.go            # JWT + bcrypt functions
├── Dockerfile         # API multistage Dockerfile
├── docker-compose.yml # All services
├── go.mod
├── go.sum
└── consumer/
    ├── main.go        # Consumer / worker
    └── Dockerfile     # Consumer multistage Dockerfile

🚀 Getting Started
Prerequisites
Docker installed and running
Run everything with one command:
docker compose up --build

This starts all 5 services:
Postgres on port 5432
Redis on port 6379
RabbitMQ on port 5672
API on port 8080
Consumer (background)


📡 API Endpoints
Register
POST /register

Request body:
{
    "email": "lipika@gmail.com",
    "password": "abc123"
}

Response:
{
    "message": "user created"
}

Login
POST /login

Request body:
{
    "email": "lipika@gmail.com",
    "password": "abc123"
}

Response:
{
    "token": "eyJhbGciOiJIUzI1NiJ9..."
}

Shorten a URL (protected 🔐)
POST /shorten
Authorization: Bearer <token>

Request body:
{
    "url": "https://google.com"
}

Response:
{
    "short_url": "http://localhost:8080/abc123"
}

Redirect to original URL
GET /:code

Redirects to the original URL and tracks the click.




⚙️ Environment Variables
Variable                            Default                                                 Description
DATABASE_URL                        postgres://postgres:pass@localhost:5432/postgres        Postgres connection string
REDIS_URL                           localhost:6379                                          Redis address
RABBITMQ_URL                        amqp://guest:guest@localhost:5672/                     RabbitMQ connection string
JWT_SECRET                          my-secret-key                                           Secret key for signing JWT tokens

🔐 Authentication Flow
1. POST /register → password hashed with bcrypt → saved to Postgres
2. POST /login → password verified → JWT token returned
3. POST /shorten → token validated → URL created

🛡️ Rate Limiting
Every IP address is limited to 10 requests per minute.
Request 1-10  → 200 OK ✅
Request 11+   → 429 Too Many Requests ❌
After 1 min   → counter resets ✅

🔄 How It Works
Shorten URL flow:
JWT middleware validates token
POST /shorten receives a long URL
Generates a random 6-character code
Saves to Postgres
Caches in Redis immediately
Returns the short URL
Redirect flow:
GET /:code receives a short code
Checks Redis first (cache hit → instant redirect ⚡)
On cache miss → queries Postgres → saves to Redis
Publishes click event to RabbitMQ
Redirects user to original URL
Consumer flow:
Listens to RabbitMQ clicks queue
Picks up click events
Increments click count in Postgres



🗄️ Database Schema
CREATE TABLE urls (
    code         TEXT PRIMARY KEY,
    original_url TEXT,
    clicks       INT DEFAULT 0
);

CREATE TABLE users (
    id       SERIAL PRIMARY KEY,
    email    TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

🧪 Testing
Use Postman or curl:
# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "lipika@gmail.com", "password": "abc123"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "lipika@gmail.com", "password": "abc123"}'

# Shorten a URL (use token from login)
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"url": "https://google.com"}'

# Use the short URL
curl -L http://localhost:8080/abc123


