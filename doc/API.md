# Quest Manager API Documentation

## 📡 API Overview

Quest Manager provides RESTful API for managing quests with location-based features.

**Base URL:** `http://localhost:8080`  
**API Version:** `v1`  
**Protocol:** HTTP/REST  
**Format:** JSON  
**Authentication:** JWT Bearer Token (required for all endpoints)

---

## 🔐 Authentication

All API endpoints require JWT authentication via Bearer token.

### Header Format
```http
Authorization: Bearer <your-jwt-token>
```

### Authentication Flow
1. Client obtains JWT token from Quest Auth service
2. Client includes token in `Authorization` header
3. API validates token via gRPC call to Auth service
4. User ID is extracted from token and injected into request context
5. Handlers use authenticated user ID for operations

### Error Responses

#### 401 Unauthorized - Missing Token
```json
{
  "type": "about:blank",
  "title": "Authentication Failed",
  "status": 401,
  "detail": "Invalid or malformed authentication token"
}
```

#### 401 Unauthorized - Expired Token
```json
{
  "type": "about:blank",
  "title": "Token Expired",
  "status": 401,
  "detail": "JWT token has expired, please refresh your token"
}
```

---

## 📋 Endpoints

### Quest Creation

#### `POST /api/v1/quests`
Create a new quest.

**Authentication:** Required  
**User ID Source:** JWT token (creator)

**Request Body:**
```json
{
  "title": "Find the treasure",
  "description": "Search for hidden treasure in the park",
  "difficulty": "medium",
  "reward": 3,
  "duration_minutes": 120,
  "target_location": {
    "latitude": 55.7558,
    "longitude": 37.6173,
    "address": "Red Square, Moscow"
  },
  "execution_location": {
    "latitude": 55.7558,
    "longitude": 37.6173,
    "address": "Red Square, Moscow"
  },
  "equipment": ["map", "compass"],
  "skills": ["navigation", "orienteering"]
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Find the treasure",
  "description": "Search for hidden treasure in the park",
  "difficulty": "medium",
  "reward": 3,
  "duration_minutes": 120,
  "status": "created",
  "creator": "user-id-from-token",
  "assignee": null,
  "created_at": "2025-10-09T10:30:00Z",
  "updated_at": "2025-10-09T10:30:00Z",
  "target_location": {...},
  "execution_location": {...},
  "equipment": ["map", "compass"],
  "skills": ["navigation", "orienteering"]
}
```

---

### Quest Retrieval

#### `GET /api/v1/quests`
Get list of all quests with optional status filter.

**Authentication:** Required

**Query Parameters:**
- `status` (optional): Filter by status (`created`, `posted`, `assigned`, `in_progress`, `declined`, `completed`)

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Find the treasure",
    "status": "posted",
    ...
  }
]
```

---

#### `GET /api/v1/quests/{quest_id}`
Get quest details by ID.

**Authentication:** Required

**Path Parameters:**
- `quest_id`: UUID of the quest

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Find the treasure",
  ...
}
```

**Error Responses:**
- `404 Not Found` - Quest doesn't exist

---

#### `GET /api/v1/quests/assigned`
Get quests assigned to authenticated user.

**Authentication:** Required  
**User ID Source:** JWT token (automatic)

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "assignee": "user-id-from-token",
    "status": "assigned",
    ...
  }
]
```

**Note:** Returns only quests assigned to the authenticated user.

---

#### `GET /api/v1/quests/search-radius`
Search quests within geographic radius.

**Authentication:** Required

**Query Parameters:**
- `lat` (required): Center latitude (-90 to 90)
- `lon` (required): Center longitude (-180 to 180)
- `radius_km` (required): Search radius in kilometers (0.1 to 20000)

**Example:**
```http
GET /api/v1/quests/search-radius?lat=55.7558&lon=37.6173&radius_km=10
```

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "target_location": {
      "latitude": 55.7558,
      "longitude": 37.6173
    },
    ...
  }
]
```

---

### Quest Assignment

#### `POST /api/v1/quests/{quest_id}/assign`
Assign quest to authenticated user.

**Authentication:** Required  
**User ID Source:** JWT token (automatic)

**Path Parameters:**
- `quest_id`: UUID of the quest to assign

**Request Body:** None (user ID taken from JWT token)

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "assignee": "user-id-from-token",
  "status": "assigned"
}
```

**Error Responses:**
- `404 Not Found` - Quest doesn't exist
- `400 Bad Request` - Quest already assigned or invalid status

---

### Quest Status Management

#### `PATCH /api/v1/quests/{quest_id}/status`
Change quest status.

**Authentication:** Required

**Path Parameters:**
- `quest_id`: UUID of the quest

**Request Body:**
```json
{
  "status": "posted"
}
```

**Valid Status Transitions:**
```
created → posted, assigned
posted → created, assigned
assigned → posted, in_progress, declined
in_progress → completed, declined
declined → posted
completed → (terminal state, no transitions)
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "posted",
  "assignee": null
}
```

**Error Responses:**
- `404 Not Found` - Quest doesn't exist
- `400 Bad Request` - Invalid status or transition

---

## 🎯 Quest Status Lifecycle

```
created → posted → assigned → in_progress → completed
            ↑         ↓            ↓
            └─────  declined ──────┘
```

**Status Descriptions:**
- `created` - Quest just created, not yet published
- `posted` - Quest is public and available for assignment
- `assigned` - Quest is assigned to a user
- `in_progress` - User started working on the quest
- `declined` - User declined the quest
- `completed` - Quest successfully finished (terminal state)

---

## 🌍 Location Features

### Coordinate System
- **Latitude:** -90 to 90 degrees
- **Longitude:** -180 to 180 degrees
- **Precision:** Float32
- **Distance Calculation:** Haversine formula

### Two Location Types

**Target Location:**
- Where quest objective is located
- Used for geographic searches

**Execution Location:**
- Where quest should be performed
- May differ from target location

---

## 📊 Common Response Patterns

### Success Response (2xx)
```json
{
  "id": "uuid",
  "status": "string",
  ...
}
```

### Validation Error (400)
```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 400,
  "detail": "Invalid input: field 'difficulty' must be one of [easy, medium, hard]"
}
```

### Not Found (404)
```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "quest with id '550e8400-...' not found"
}
```

### Server Error (500)
```json
{
  "type": "about:blank",
  "title": "Internal Server Error",
  "status": 500,
  "detail": "An unexpected error occurred"
}
```

---

## 🔧 API Client Examples

### cURL Examples

#### Create Quest
```bash
curl -X POST http://localhost:8080/api/v1/quests \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Find the treasure",
    "description": "Search for hidden treasure",
    "difficulty": "medium",
    "reward": 3,
    "duration_minutes": 120,
    "target_location": {
      "latitude": 55.7558,
      "longitude": 37.6173
    },
    "execution_location": {
      "latitude": 55.7558,
      "longitude": 37.6173
    }
  }'
```

#### Assign Quest to Self
```bash
curl -X POST http://localhost:8080/api/v1/quests/{quest-id}/assign \
  -H "Authorization: Bearer <your-token>"
```

#### List My Assigned Quests
```bash
curl http://localhost:8080/api/v1/quests/assigned \
  -H "Authorization: Bearer <your-token>"
```

#### Search Quests Near Moscow
```bash
curl "http://localhost:8080/api/v1/quests/search-radius?lat=55.7558&lon=37.6173&radius_km=10" \
  -H "Authorization: Bearer <your-token>"
```

---

## 📝 Field Validations

### Quest Fields

| Field              | Type          | Constraints                      | Required |
|--------------------|---------------|----------------------------------|----------|
| title              | string        | 1-200 chars, no whitespace-only  | ✅        |
| description        | string        | 1-1000 chars, no whitespace-only | ✅        |
| difficulty         | enum          | easy, medium, hard               | ✅        |
| reward             | integer       | 1-5                              | ✅        |
| duration_minutes   | integer       | 1-10080 (1 week)                 | ✅        |
| target_location    | object        | Valid coordinates                | ✅        |
| execution_location | object        | Valid coordinates                | ✅        |
| equipment          | array[string] | Max 50 items, 1-100 chars each   | ❌        |
| skills             | array[string] | Max 50 items, 1-100 chars each   | ❌        |

### Coordinate Fields

| Field     | Type   | Constraints | Required |
|-----------|--------|-------------|----------|
| latitude  | float  | -90 to 90   | ✅        |
| longitude | float  | -180 to 180 | ✅        |
| address   | string | 1-500 chars | ❌        |

---

## 🔗 Related Documentation

- [Architecture](ARCHITECTURE.md) - System architecture and design patterns
- [OpenAPI Specification](../api/http/quests/v1/openapi.yaml) - Full API specification
- [Changelog](changelog/) - Version history and changes
- [Testing Guide](TESTING.md) - Testing strategies and patterns

---

## 💡 Best Practices

### 1. Always Use Valid UUIDs
All ID fields must be valid UUID v4 format.

### 2. Handle Status Transitions
Check valid status transitions before changing quest status.

### 3. Coordinate Validation
Ensure latitude/longitude are within valid ranges.

### 4. Token Management
- Store tokens securely
- Refresh tokens before expiration
- Handle 401 errors gracefully

### 5. Error Handling
- Parse Problem Details format
- Display user-friendly error messages
- Log errors for debugging

---

**Last Updated:** October 9, 2025  
**API Version:** 1.4.0

