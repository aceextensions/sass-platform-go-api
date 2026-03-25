# Database Migrations Guide

This document explains how to manage database schema changes in the AceExtension Go Platform.

## 🏗️ Architecture

The platform uses a **Distributed Migration System**. Instead of a single central migration folder, each module (e.g., `accounting`, `catalog`, `identity`) manages its own schema changes within its own `migrations` directory.

### How Discovery Works
The migration CLI (`api/cmd/migrate/main.go`) performs a **recursive workspace scan**:
1. It starts from the project root (detecting `go.work`).
2. It walks through all subdirectories (skipping hidden ones).
3. It identifies every folder named `migrations`.
4. It collects all `.sql` files and sorts them **lexicographically by their relative path**.
5. It tracks applied migrations in the `schema_migrations` table using the relative path (e.g., `core/migrations/001_initial_schema.sql`) as a unique identifier.

## 🚀 Commands

You can run migrations either locally or via Docker Compose.

### Apply Pending Migrations
```bash
# Locally
go run ./api/cmd/migrate/main.go up

# via Docker
docker compose exec api-go go run ./api/cmd/migrate/main.go up
```

### Check Migration Status
```bash
# Locally
go run ./api/cmd/migrate/main.go status

# via Docker
docker compose exec api-go go run ./api/cmd/migrate/main.go status
```

## ➕ Adding a New Migration

Follow these steps to add a schema change:

### 1. Identify the Module
Decide which module owns the tables you are modifying. 
- If it's a global table (users, tenants, etc.), use `core/migrations`.
- If it's specific to a feature, use that module's `migrations` folder (e.g., `quiz/migrations`).

### 2. Create the SQL File
Use a numeric prefix to ensure correct ordering.
- **Good Name**: `004_add_profile_image_to_users.sql`
- **Bad Name**: `add_image.sql`

### 3. Write Idempotent SQL
Always use `IF NOT EXISTS` for tables/columns and `CREATE OR REPLACE` for functions/triggers. This ensures the migration can safely run multiple times without failing.

Example:
```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
```

### 4. Commit and Run
Commit the new `.sql` file to Git and run the `up` command to apply it to your local database.

## 💡 Best Practices

1. **Atomic Migrations**: Each file should represent one logical change.
2. **No Data Manipulation**: Try to keep migrations focused on **schema** (DDL). Avoid complex data migrations (DML) inside SQL files if possible; use a Go-based script for complex data transformations.
3. **Rollbacks**: Currently, the tool is "up-only". If you make a mistake, create a new "fix" migration instead of deleting the old one.
4. **Build-Time Independence**: The application code should **never** require a running database to compile. All DB setup must happen at runtime via these migrations.
