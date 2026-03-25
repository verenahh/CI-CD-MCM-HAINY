# Exercise 2: Microservice Architecture, Docker & GitHub Actions

**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 24

## Learning Objectives

- Understand microservice architecture with a REST API in Go
- Containerize applications using Docker (multi-stage builds)
- Orchestrate services with Docker Compose
- Set up a basic CI pipeline with GitHub Actions

## Prerequisites

- Completed Exercise 1
- Docker Desktop installed
- Basic understanding of REST APIs

## Project Overview

The Product Catalog API has been extended with:
- **PostgreSQL storage** (`internal/store/postgres.go`) -- persistent database backend
- **Dockerfile** -- multi-stage build for minimal container image
- **docker-compose.yml** -- orchestrates API + PostgreSQL
- **GitHub Actions** (`.github/workflows/ci.yml`) -- basic CI pipeline

### Architecture

```
┌──────────────┐     ┌──────────────┐
│   Client     │────▶│   API (Go)   │
│  (curl/HTTP) │     │   Port 8080  │
└──────────────┘     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  PostgreSQL  │
                     │  Port 5432   │
                     └──────────────┘
```

### Local Development

```bash
# Run with in-memory store (no Docker needed)
go run ./cmd/api

# Run with Docker Compose (API + PostgreSQL)
docker compose up --build

# Test the API
curl http://localhost:8080/health
curl http://localhost:8080/products
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99}'
```

Each exercise branch contains a detailed `README.md` with instructions.

## Authors
- Prof. M. Kurz
- Student Name
