# Social Login Setup Guide

This guide explains how to configure and enable Social Login (OAuth2/OIDC) for the platform using the `sociallogin-pkg`.

## Supported Providers

- **Google** (OpenID Connect)
- **GitHub** (OAuth2)

## Configuration

Social login is enabled by providing the necessary credentials in your `.env` file. If credentials for a provider are missing, the corresponding authentication route will return an error.

### Required Environment Variables

Add the following to your `.env` file:

```env
# Social Login Credentials
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# Session Secret (Required by Gorilla Sessions / Goth)
SESSION_SECRET=your-random-session-secret
```

### Redirect / Callback URLs

You must register the following callback URLs in your provider's developer console (Google Cloud Console / GitHub Settings):

| Provider | Callback URL |
| :--- | :--- |
| **Google** | `https://api.practixa.com/api/auth/google/callback` |
| **GitHub** | `https://api.practixa.com/api/auth/github/callback` |

> [!IMPORTANT]
> Ensure that the domain matches your `APP_BASE_URL`. For local development, use `http://localhost:3000/api/auth/[provider]/callback`.

## Provider Setup Instructions

### 1. Google Cloud Console
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project or select an existing one.
3. Search for **APIs & Services** > **Credentials**.
4. Click **Create Credentials** > **OAuth client ID**.
5. Select **Web application** as the application type.
6. Add the **Authorized redirect URIs** listed above.
7. Copy the **Client ID** and **Client Secret** to your `.env`.

### 2. GitHub Developer Settings
1. Log in to GitHub and go to **Settings** > **Developer settings** > **OAuth Apps**.
2. Click **New OAuth App**.
3. Set the **Homepage URL** (e.g., `https://practixa.com`).
4. Set the **Authorization callback URL** listed above.
5. Click **Register application**.
6. Generate a **Client Secret**.
7. Copy the **Client ID** and **Client Secret** to your `.env`.

## User Provisioning Logic

When a user logs in via a social provider:
1. The system checks if a user exists with the verified email address.
2. If the user exists, the social account is linked and the session is started.
3. If the user does not exist, a new **Individual** account is automatically provisioned, and a personal tenant is created for the user.
