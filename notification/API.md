# Notification Module API Documentation

Base URL: `/api/v1`

## 1. Notifications

### Send Notification
**POST** `/notifications/send`

Trigger a new notification (Internal/Admin use).

**Request Body:**
```json
{
  "userId": "uuid",
  "type": "EMAIL", // EMAIL, SMS, IN_APP
  "subject": "Welcome!",
  "message": "Thanks for joining...",
  "priority": "HIGH"
}
```

### List Notifications
**GET** `/notifications`

Retrieve notifications for the current user.

**Query Params:**
- `unreadOnly`: boolean

### Mark as Read
**PUT** `/notifications/{id}/read`

Mark a specific notification as read.
