# Social Login & Authentication Guide

This guide explains the full authentication architecture for the platform, including Social Login (OAuth2/OIDC), JWT token management, and API endpoint reference.

## Architecture Overview

| Component | Technology | Notes |
| :--- | :--- | :--- |
| **Token Format** | JWT (HS256) | `golang-jwt/jwt/v5` |
| **Password Hashing** | Argon2id | CPU/memory-hard KDF |
| **Session Storage** | Database (`sessions` table) | Refresh tokens stored server-side |
| **Social Login** | Gothic/Goth | OpenID Connect & OAuth2 |
| **Middleware** | Custom Echo Middleware | Bearer token extraction |

> [!IMPORTANT]
> The platform uses **JWT (JSON Web Tokens)** — not PASETO. All tokens are signed with HMAC-SHA256 using the `JWT_SECRET` environment variable.

---

## API Endpoints Reference

All auth endpoints live under `/api/v1/auth`.

### Public Endpoints (No Token Required)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Register a new tenant (business) |
| `POST` | `/api/v1/auth/register-individual` | Register a personal account |
| `POST` | `/api/v1/auth/verify-otp` | Verify registration OTP |
| `POST` | `/api/v1/auth/login` | Login with email/phone + password |
| `POST` | `/api/v1/auth/refresh` | Refresh expired access token |
| `POST` | `/api/v1/auth/forgot-password` | Request password reset OTP |
| `POST` | `/api/v1/auth/reset-password` | Reset password with OTP |
| `GET` | `/api/v1/auth/:provider` | Begin social login (Google/GitHub) |
| `GET` | `/api/v1/auth/:provider/callback` | Social login callback |

### Protected Endpoints (Bearer Token Required)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/auth/me` | Get current authenticated user profile |
| `POST` | `/api/v1/auth/logout` | Invalidate session |
| `POST` | `/api/v1/auth/change-password` | Change authenticated user's password |
| `POST` | `/api/v1/auth/impersonate/:tenantId` | Login as tenant owner (Super Admin only) |

---

## Token Structure

### Access Token (1 hour TTL)

```json
{
  "userId": "a3da1839-1bfb-4dc5-85f8-62822e337e7b",
  "tenantId": "54eeb799-2c23-4119-97f4-2f1b06247e3f",
  "role": "owner",
  "exp": 1774701909,
  "iat": 1774698309
}
```

### Refresh Token (7 day TTL)

```json
{
  "userId": "a3da1839-1bfb-4dc5-85f8-62822e337e7b",
  "tenantId": "54eeb799-2c23-4119-97f4-2f1b06247e3f",
  "role": "owner",
  "exp": 1775303109,
  "iat": 1774698309
}
```

### Using the Refresh Token

When the access token expires, call:

```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "<your-refresh-token>"
}
```

**Response:**
```json
{
  "accessToken": "eyJhbG...",
  "refreshToken": "eyJhbG...",
  "user": {
    "id": "a3da1839-...",
    "name": "Durga Gupta",
    "email": "durga@example.com",
    "role": "owner",
    "tenantId": "54eeb799-..."
  }
}
```

> [!NOTE]
> The refresh endpoint performs **token rotation**: the old refresh token is invalidated and a new pair (access + refresh) is returned. This prevents replay attacks.

---

## Social Login Flow

### Supported Providers

- **Google** (OpenID Connect)
- **GitHub** (OAuth2)

### Flow Sequence

```
Browser → GET /api/v1/auth/google
       → Google Consent Screen
       → GET /api/v1/auth/google/callback
       → Server creates/links user + generates JWT
       → 307 Redirect to APP_BASE_URL/login/success?accessToken=...&refreshToken=...
       → Frontend stores tokens
       → GET /api/v1/auth/me (with Bearer token header)
```

### Configuration

Social login is enabled by providing the necessary credentials in your `.env` file. If credentials for a provider are missing, the corresponding authentication route will return an error.

#### Required Environment Variables

Add the following to your `.env` file:

```env
# Social Login Credentials
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# Session Secret (Required by Gorilla Sessions / Goth)
SESSION_SECRET=your-random-session-secret

# Base URLs (used for callback construction and redirects)
API_BASE_URL=https://api.example.com
APP_BASE_URL=https://dashboard.example.com
```

> [!WARNING]
> Never hardcode domain names in source code. All URLs are constructed from `API_BASE_URL` and `APP_BASE_URL` environment variables.

#### Redirect / Callback URLs

You must register the following callback URLs in your provider's developer console (Google Cloud Console / GitHub Settings):

| Provider | Callback URL |
| :--- | :--- |
| **Google** | `${API_BASE_URL}/api/v1/auth/google/callback` |
| **GitHub** | `${API_BASE_URL}/api/v1/auth/github/callback` |

> [!IMPORTANT]
> Ensure that the domain matches your `API_BASE_URL`. For local development, use `http://localhost:4000/api/v1/auth/[provider]/callback`.

### Provider Setup Instructions

#### 1. Google Cloud Console
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project or select an existing one.
3. Search for **APIs & Services** > **Credentials**.
4. Click **Create Credentials** > **OAuth client ID**.
5. Select **Web application** as the application type.
6. Add the **Authorized redirect URIs** listed above.
7. Copy the **Client ID** and **Client Secret** to your `.env`.

#### 2. GitHub Developer Settings
1. Log in to GitHub and go to **Settings** > **Developer settings** > **OAuth Apps**.
2. Click **New OAuth App**.
3. Set the **Homepage URL** to your `APP_BASE_URL` value.
4. Set the **Authorization callback URL** listed above.
5. Click **Register application**.
6. Generate a **Client Secret**.
7. Copy the **Client ID** and **Client Secret** to your `.env`.

### User Provisioning Logic

When a user logs in via a social provider:
1. The system checks if a user exists with the verified email address.
2. If the user exists, the social account is linked and a session is started.
3. If the user does not exist, a new **Individual** account is automatically provisioned, a personal tenant is created, and the user is auto-verified.

---

## Authentication Middleware

The `JWTMiddleware` extracts and validates Bearer tokens from the `Authorization` header.

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Behavior:**
- Extracts token from `Authorization: Bearer <token>` header
- Validates signature using `JWT_SECRET`
- Verifies expiration (`exp` claim)
- Injects `AuthUser{UserID, TenantID, Role}` into Echo context
- Returns `401 Unauthorized` if token is missing, invalid, or expired

### Role-Based Access Control

Use `RequireRole(roles...)` middleware after `JWTMiddleware` to restrict endpoints:

```go
protectedGroup.Use(middleware.JWTMiddleware)
protectedGroup.POST("/admin-action", handler, middleware.RequireRole("owner", "admin"))
```

---

## Current Feature Status

### ✅ Implemented
- JWT token generation (Access + Refresh)
- Token refresh with rotation
- Social login (Google, GitHub)
- Argon2id password hashing
- Session management (DB-backed)
- JWTMiddleware with Bearer extraction
- Role-based access control middleware
- OTP verification flow
- Password reset flow

### ❌ Not Yet Implemented

#### API Token Management
- No `api_tokens` database table
- No long-lived API token generation for programmatic access
- No API token validation middleware
- **Planned**: `POST /api/v1/auth/api-tokens` to create, `DELETE` to revoke

#### Workspace Switching Endpoint
- No `POST /api/v1/workspace/:id/select` endpoint
- No JWT claim update to switch active tenant context
- **Planned**: Allows users with multiple memberships to switch workspace context by issuing new tokens with updated `tenantId` claim
