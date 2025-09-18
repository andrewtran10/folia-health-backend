# Folia Health Backend API Documentation

## Overview

This API provides endpoints for managing users and reminders in the Folia Health application. The API uses Laravel Sanctum for authentication and supports CRUD operations for both users and reminders with recurring scheduling capabilities.

### Authentication

The API uses Laravel Sanctum token-based authentication. Most endpoints require authentication via the `Authorization: Bearer {token}` header.

#### Getting an API Token

When creating a user, an API token is automatically generated and returned in the response. Use this token for subsequent authenticated requests.

---

## Endpoints

### Users

#### Create User

**POST** `/api/users`

Creates a new user account and returns an API token for authentication.

**Authentication:** Not required

**Request Body:**

```json
{
    "name": "string (required)",
    "email": "string (required)",
    "password": "string (required)"
}
```

**Response (201):**

```json
{
    "data": {
        "id": "uuid",
        "name": "string",
        "email": "string",
        "created_at": "datetime",
        "updated_at": "datetime",
        "token": "string"
    }
}
```

**Example:**

```bash
curl -X POST /api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Andrew Tran",
    "email": "andrew_tran@realemail.com",
    "password": "password123"
  }'
```

#### Get User

**GET** `/api/users/{id}`

Retrieves a specific user by their ID.

**Authentication:** Required

**Response (200):**

```json
{
    "data": {
        "id": "uuid",
        "name": "string",
        "email": "string",
        "created_at": "datetime",
        "updated_at": "datetime"
    }
}
```

**Example:**

```bash
curl -X GET /api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer {token}"
```

---

### Reminders

#### Create Reminder

**POST** `/api/reminders`

Creates a new reminder with recurrence rules.

**Authentication:** Required

**Request Body:**

```json
{
    "rrule": "string (required) - RRULE formatted string (RFC 5545)",
    "description": "string (required) - Reminder description",
    "start_at": "string (required) - datetime"
}
```

**Response (200):**

```json
{
    "data": {
        "id": "uuid",
        "rrule": "string",
        "rrule_human": "string - Human readable recurrence description",
        "description": "string",
        "start_at": "datetime"
    }
}
```

**RRULE Format Examples:**

-   Daily: `FREQ=DAILY;INTERVAL=1`
-   Weekly on Monday and Wednesday: `FREQ=WEEKLY;BYDAY=MO,WE`
-   Monthly on the 15th: `FREQ=MONTHLY;BYMONTHDAY=15`
-   Every 2 weeks: `FREQ=WEEKLY;INTERVAL=2`

**Example:**

```bash
curl -X POST /api/reminders \
  -H "Authorization: Bearer {api-token}" \
  -H "Content-Type: application/json" \
  -d '{
    "rrule": "FREQ=DAILY;INTERVAL=1",
    "description": "Take medication",
    "start_at": "2025-09-17T09:00:00Z"
  }'
```

#### Get Reminders

**GET** `/api/reminders`

Retrieves reminders for the authenticated user with optional filtering.

**Authentication:** Required

**Query Parameters:**

-   `search` (optional): Search term to filter by description (case-insensitive)
-   `start_date` (optional): Start date for date range filter (requires end_date)
-   `end_date` (optional): End date for date range filter (requires start_date)

If filtering by range both start_date and end_date are required.

**Response (200):**

```json
{
    "data": [
        {
            "id": "uuid",
            "rrule": "string",
            "rrule_human": "string",
            "description": "string",
            "start_at": "datetime",
            "occurrences": ["datetime"] // Only included when start_date and end_date are provided
        }
    ]
}
```

**Examples:**

Get all reminders:

```bash
curl -X GET /api/reminders \
  -H "Authorization: Bearer {api-token}"
```

Search reminders:

```bash
curl -X GET "/api/reminders?search=medication" \
  -H "Authorization: Bearer {api-token}"
```

Get reminders with occurrences in date range:

```bash
curl -X GET "/api/reminders?start_date=2025-09-17&end_date=2025-09-24" \
  -H "Authorization: Bearer {api-token}"
```

#### Update Reminder

**PATCH** `/api/reminders/{reminder}`

Updates an existing reminder.

**Authentication:** Required

**Request Body:**

```json
{
    "rrule": "string (optional) - RRULE formatted string (RFC 5545)",
    "description": "string (optional) - Reminder description",
    "start_at": "string (optional) - datetime format"
}
```

**Response (200):**

```json
{
    "data": {
        "id": "uuid",
        "rrule": "string",
        "rrule_human": "string",
        "description": "string",
        "start_at": "datetime"
    }
}
```

**Example:**

```bash
curl -X PATCH /api/reminders/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer {api-token}" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Take medication with food"
  }'
```

#### Delete Reminder

**DELETE** `/api/reminders/{reminder}`

Deletes a specific reminder.

**Authentication:** Required

**Response (204):** No content

**Example:**

```bash
curl -X DELETE /api/reminders/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer {api-token}"
```

---

## RRULE Reference

The API uses the RRULE standard (RFC 5545) for defining recurring reminders. Here are common patterns:

### Frequency Types

-   `FREQ=SECONDLY`
-   `FREQ=MINUTELY`
-   `FREQ=HOURLY`
-   `FREQ=DAILY`
-   `FREQ=WEEKLY`
-   `FREQ=MONTHLY`
-   `FREQ=YEARLY`

### Common Parameters

-   `INTERVAL=n` - Repeat every n intervals
-   `COUNT=n` - Stop after n occurrences
-   `UNTIL=datetime` - Stop after specified date
-   `BYDAY=MO,TU,WE,TH,FR,SA,SU` - Specific days of week
-   `BYMONTHDAY=1,15` - Specific days of month

### Examples

-   Every day: `FREQ=DAILY`
-   Every other day: `FREQ=DAILY;INTERVAL=2`
-   Every weekday: `FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR`
-   First Monday of every month: `FREQ=MONTHLY;BYDAY=1MO`
-   Every 3 months on the 15th: `FREQ=MONTHLY;INTERVAL=3;BYMONTHDAY=15`

---
