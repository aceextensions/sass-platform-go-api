# Subscription Module API Documentation

Base URL: `/api/v1`

## 1. Plans

### Create Plan
**POST** `/plans`

Create a new subscription plan (Admin only).

**Request Body:**
```json
{
  "name": "Gold Tier",
  "description": "Premium features",
  "price": 99.99,
  "billingCycle": "MONTHLY", // MONTHLY, YEARLY
  "features": ["Feature A", "Feature B"]
}
```

### List Plans
**GET** `/plans`

Retrieve all available subscription plans.

---

## 2. Subscriptions

### Subscribe
**POST** `/subscriptions/subscribe`

Subscribe the current tenant to a plan.

**Request Body:**
```json
{
  "planId": "uuid",
  "password": "user-password-for-verification"
}
```

### Get Current Subscription
**GET** `/subscriptions/current`

Retrieve details of the active subscription for the current tenant.
