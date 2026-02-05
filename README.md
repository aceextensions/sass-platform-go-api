# AceExtension Go API 🚀

This repository contains the high-performance Go (Golang) implementation of the AceExtension API, migrating from the original Bun/Elysia backend for enhanced scalability, enterprise-grade stability, and type safety.

## 🏗️ Architecture: Modular Monorepo

We use **Go Workspaces (`go.work`)** to manage multiple independent modules within a single repository. This allows for clean separation of concerns while sharing local dependencies.

| Module | Purpose |
| :--- | :--- |
| **`api`** | Application entry point, routing, and HTTP handlers. |
| **`identity`** | Auth & User management (Repositories, Services, DTOs). |
| **`core`** | Infrastructure layer: Configuration (Viper), DB Pool (pgx), Logger (Zap). |
| **`common`** | Shared utilities, constants, and global errors. |

## 🛠️ Technology Stack

- **Framework**: [Echo v4](https://echo.labstack.com/) (High performance, minimalist)
- **Database**: [pgx v5](https://github.com/jackc/pgx) (Native PostgreSQL driver with connection pooling)
- **Configuration**: [Viper](https://github.com/spf13/viper) (Environment-aware config management)
- **Logging**: [Uber-Zap](https://github.com/uber-go/zap) (Blazing fast structured logging)
- **Documentation**: [Swaggo](https://github.com/swaggo/swag) (Swagger 2.0/OpenAPI)
- **Development**: [Air](https://github.com/air-verse/air) (Live reloading)

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (Optional for local development)

### Running with Docker (Recommended)
The Go API is integrated into the main `docker-compose.yml` at the project root.

```bash
# From the project root
docker compose up --build api-go
```

The API will be available at: **[http://localhost:4000](http://localhost:4000)**

### API Documentation (Swagger)
Interactive API docs are automatically generated and served at:
**[http://localhost:4000/swagger/index.html](http://localhost:4000/swagger/index.html)**

> [!NOTE]
> Documentation is regenerated on the fly during development whenever code changes are detected by Air.

## 🛠️ Development Workflow

### Hot Reloading
We use **Air** for instant feedback. It watches all `.go` files in the workspace and rebuilds the binary inside the container on save.

### Adding New Dependencies
Since we use workspaces, you should add dependencies to the specific module's `go.mod`:

```bash
cd identity
go get github.com/example/package
go mod tidy
```

Then sync the workspace from the root:
```bash
go work sync
```

## 🔒 Security & Auth
- **JWT**: Stateless authentication using RSA/HMAC.
- **RBAC**: Role-based access control middleware implemented in the `identity` module.
- **Migrations**: Database schema is managed via the main project's Drizzle migrations.

## 📈 Roadmap
- [x] Core Infrastructure (Config, DB, Logger)
- [x] Auth & Identity Porting
- [x] Swagger Integration
- [ ] Asynchronous Audit Logging (Goroutines + Channels)
- [ ] Production Multi-stage Build Optimization
- [ ] Direct sqlc integration for type-safe SQL
