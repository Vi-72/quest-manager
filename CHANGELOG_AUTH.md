# Authentication Integration - Changelog

## 🔐 Version 1.4.0 - JWT Authentication

### ✨ New Features

#### **JWT Authentication Middleware**
- Added global JWT authentication for all API endpoints
- Integration with Quest Auth service via gRPC
- Bearer token validation in `Authorization` header
- User ID extraction and injection into request context

#### **OpenAPI Specification Updates**
- Added `securitySchemes.bearerAuth` definition
- Added global `security: bearerAuth` requirement
- Added `401 Unauthorized` responses to all 7 endpoints
- Updated API version: `1.3.4` → `1.4.0`
- Updated description: "All endpoints require JWT authentication"

### 🔧 Technical Changes

#### **New Components**
1. **Auth Middleware** (`internal/adapters/in/http/middleware/authentication.go`)
   - Bearer token extraction from headers
   - gRPC call to Auth service
   - User ID context injection
   - Structured error handling (Problem Details RFC 7807)

2. **Auth Client** (`internal/adapters/out/client/auth/client.go`)
   - gRPC client wrapper for Quest Auth service
   - Token validation
   - User ID extraction
   - Error handling (token expired, invalid token)

3. **Auth Port** (`internal/core/ports/auth_client.go`)
   - Interface for authentication
   - DDD/Hexagonal architecture compliance

#### **Updated Components**
1. **Composition Root** (`cmd/composition_root.go`)
   - Auth gRPC connection setup
   - Auth client initialization
   - Middleware registration
   - Optional auth (disabled if `AUTH_GRPC` not set)

2. **Configuration** (`cmd/config.go`)
   - Added `AuthGRPC` field for gRPC address

3. **Error Handling** (`internal/adapters/in/http/errors/problem_details.go`)
   - Added `NewProblem()` helper for custom errors

### 📝 Configuration

#### **New Environment Variable**
```bash
AUTH_GRPC=localhost:50051  # gRPC address of Quest Auth service
```

If `AUTH_GRPC` is not set, authentication is disabled (for local development).

### 🔒 Security

#### **Protected Endpoints** (require JWT)
- `POST /api/v1/quests`
- `GET /api/v1/quests`
- `GET /api/v1/quests/{quest_id}`
- `PATCH /api/v1/quests/{quest_id}/status`
- `POST /api/v1/quests/{quest_id}/assign`
- `GET /api/v1/quests/assigned`
- `GET /api/v1/quests/search-radius`

#### **Public Endpoints** (no auth)
- `GET /health`
- `GET /docs`
- `GET /openapi.json`

### 📊 Error Responses

#### **401 Unauthorized**
```json
{
  "type": "about:blank",
  "title": "Token Expired",
  "status": 401,
  "detail": "JWT token has expired, please refresh your token"
}
```

```json
{
  "type": "about:blank",
  "title": "Authentication Failed",
  "status": 401,
  "detail": "Invalid or malformed authentication token"
}
```

### 🔄 Migration Guide

#### **For Clients**
All API requests must now include JWT token:
```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
     http://localhost:8080/api/v1/quests
```

#### **For Developers**
1. Set `AUTH_GRPC` in `.env`:
   ```bash
   AUTH_GRPC=localhost:50051
   ```

2. Ensure Quest Auth service is running:
   ```bash
   # Auth service must be available at configured address
   ```

3. For local development without auth:
   ```bash
   # Comment out or remove AUTH_GRPC from .env
   # AUTH_GRPC=
   ```

### 📚 Dependencies

#### **New Dependencies**
- `google.golang.org/grpc` - gRPC framework
- `github.com/Vi-72/quest-auth/api/grpc/sdk/go/auth/v1` - Quest Auth SDK

### 🧪 Testing Notes

**Integration tests need updating** to include JWT tokens in requests.

Current implementation uses `request.Body.UserId` for quest assignment. 
Future enhancement: use authenticated user ID from context for improved security.

### 🎯 Future Improvements

1. **Role-Based Access Control (RBAC)**
   - Define user roles (admin, user)
   - Implement permission checks
   - Admin can assign quests to others, users only to themselves

2. **Security Enhancements**
   - Use authenticated user from context instead of request body
   - Add authorization checks (can user perform this action?)
   - Audit logging for security events

3. **Token Management**
   - Token refresh mechanism
   - Token revocation support
   - Multiple token types (access, refresh)

### ✅ Checklist

- [x] OpenAPI spec updated with security definitions
- [x] Authentication middleware implemented
- [x] Auth gRPC client implemented
- [x] Composition root updated
- [x] Configuration updated
- [x] Documentation updated (README, config.example)
- [x] Error handling improved (Problem Details)
- [x] Code quality maintained (consistent naming, English comments)
- [x] Project builds successfully
- [ ] Integration tests updated (TODO)
- [ ] Manual testing with real Auth service (TODO)

### 🔗 Related

- OpenAPI Spec: `api/http/quests/v1/openapi.yaml`
- Auth Middleware: `internal/adapters/in/http/middleware/authentication.go`
- Auth Client: `internal/adapters/out/client/auth/client.go`
- Quest Auth Service: https://github.com/Vi-72/quest-auth

---

**Breaking Change:** All API endpoints now require authentication. Clients must update to include JWT tokens.
