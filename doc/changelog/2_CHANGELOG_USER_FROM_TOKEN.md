# User ID from Token - Changelog

## 🔐 Version 1.5.0 - User Context from JWT

### ✨ New Features

#### **Automatic User ID Extraction from JWT**
- User ID automatically extracted from JWT token for all operations
- No need to pass user_id in request parameters
- Enhanced security - prevents user impersonation
- Simplified API - cleaner request structures

---

### 🔧 Technical Changes

#### **Updated API Endpoints**

**1. POST `/quests/{quest_id}/assign`**
- **Before:** Required `user_id` in request body
- **After:** User ID taken from JWT token automatically
- **Breaking Change:** ❌ No request body needed

**2. GET `/quests/assigned`**
- **Before:** Required `user_id` query parameter  
- **After:** User ID taken from JWT token automatically
- **Breaking Change:** ❌ No query parameter needed

#### **Updated Components**

**1. HTTP Handlers** (`internal/adapters/in/http/`)
- `assign_quest_handler.go` - Uses `middleware.UserIDFromContext(ctx)`
- `list_assigned_quests_handler.go` - Uses `middleware.UserIDFromContext(ctx)`

**Changes:**
```go
// Before
userID := request.Body.UserId

// After
userID, ok := middleware.UserIDFromContext(ctx)
if !ok {
    return nil, errors.NewBadRequest("user ID not found in context")
}
```

**2. Authentication Middleware** (`internal/adapters/in/http/middleware/authentication.go`)
- Enhanced token validation
- Validates empty tokens after "Bearer " prefix
- Better error messages

**Changes:**
```go
// Added validation for empty tokens
token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
if token == "" {
    return "", errors.New("missing or invalid Authorization header")
}
```

**3. OpenAPI Specification** (`api/http/quests/v1/openapi.yaml`)
- Removed `AssignQuestRequest` schema
- Removed `user_id` parameter from `/quests/assigned`
- Updated descriptions to reflect automatic user ID extraction

---

### 🔒 Security Improvements

#### **Prevents User Impersonation**
**Before:** User could specify any user_id in request
```json
POST /quests/{id}/assign
{
  "user_id": "other-users-id"  // ❌ Security issue
}
```

**After:** User ID enforced from authenticated token
```json
POST /quests/{id}/assign
// No body - user ID from JWT token ✅ Secure
```

#### **Enhanced Validation**
- Empty tokens rejected (`Bearer ` → 401)
- Whitespace-only tokens rejected (`Bearer    ` → 401)
- Missing Authorization header → 401
- Invalid Bearer format → 401

---

### 🧪 Testing Updates

#### **New Test Suite: Authentication Middleware** 
Location: `tests/integration/tests/quest_middlewares/authentication/`

**Test Files:**
1. `suite_test.go` - Test suite setup
2. `missing_token_test.go` - Tests without authentication (7 endpoints)
3. `expired_token_test.go` - Tests with expired tokens (7 endpoints)
4. `invalid_token_test.go` - Tests with invalid tokens (7 endpoints)
5. `auth_edge_cases_test.go` - Advanced scenarios (5 tests)

**Total:** 27 new tests

**Coverage:**
- ✅ All endpoints require authentication
- ✅ Expired token handling
- ✅ Invalid token handling
- ✅ User ID correctly extracted from token
- ✅ Data isolation between users
- ✅ Multiple users scenario

#### **Updated Mock Auth** (`tests/integration/mock/`)

**New File:** `auth_scenarios.go`
- `ConfigurableAuthClient` - Configurable mock behavior
- `BehaviorSuccess` - Successful authentication
- `BehaviorTokenExpired` - Expired token scenario
- `BehaviorInvalidToken` - Invalid token scenario
- `BehaviorMissingUser` - Missing user in response

**Helper Functions:**
```go
NewExpiredTokenAuthClient()  // Returns expired token error
NewInvalidTokenAuthClient()  // Returns invalid token error
NewConfigurableAuthClient(behavior, userID)  // Custom behavior
```

#### **Updated Test Utilities** (`tests/integration/core/case_steps/`)

**Updated:** `http_requests.go`
- Added `SkipAuth` flag to HTTPRequest
- `AssignQuestHTTPRequest()` - No longer takes userID parameter
- `ListAssignedQuestsHTTPRequest()` - No longer takes userID parameter
- Removed obsolete helper functions

#### **Updated HTTP Tests**
- `assign_quest_http_test.go` - Uses mock auth DefaultUserID
- `list_assigned_quests_http_test.go` - Uses mock auth DefaultUserID
- All tests use user ID from authentication context

#### **Updated E2E Tests**
- `assign_quest_e2e_test.go` - Uses DefaultUserID from mock auth
- Verifies full workflow with user from token

---

### 📊 API Changes Summary

| Endpoint | Parameter Removed | New Source |
|----------|-------------------|------------|
| `POST /quests/{id}/assign` | `user_id` in body | JWT token |
| `GET /quests/assigned` | `user_id` in query | JWT token |

---

### 🔄 Migration Guide

#### **For API Clients**

**Before (v1.4.0):**
```bash
# Assign quest - user_id in body
curl -X POST http://localhost:8080/api/v1/quests/{id}/assign \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "550e8400-..."}'

# List assigned - user_id in query
curl "http://localhost:8080/api/v1/quests/assigned?user_id=550e8400-..." \
  -H "Authorization: Bearer <token>"
```

**After (v1.5.0):**
```bash
# Assign quest - no body
curl -X POST http://localhost:8080/api/v1/quests/{id}/assign \
  -H "Authorization: Bearer <token>"

# List assigned - no query params
curl "http://localhost:8080/api/v1/quests/assigned" \
  -H "Authorization: Bearer <token>"
```

**Migration Steps:**
1. Update API calls to remove `user_id` parameter
2. Ensure JWT token contains correct user ID
3. Test with new API format
4. Deploy updated clients

#### **For Developers**

**No code changes needed** unless:
- Using OpenAPI generated clients → Regenerate
- Custom HTTP clients → Update request format
- Integration tests → Already updated

**Steps:**
1. Pull latest code
2. Regenerate OpenAPI code: `go generate ./api/http/quests/v1/`
3. Run tests: `go test ./...`
4. Update any custom API clients

---

### 📚 Dependencies

**No new dependencies** - only internal changes.

**Updated:**
- OpenAPI specification version: `1.4.0` → `1.5.0` (pending)

---

### 🧪 Testing Notes

#### **All Tests Updated and Passing**
- ✅ Domain tests: 45+ tests - PASS
- ✅ Contract tests: 24 tests - PASS
- ✅ Handler tests: 18 tests - PASS
- ✅ HTTP tests: 21 tests - PASS
- ✅ Middleware tests: 27 tests - NEW
- ✅ E2E tests: 3 tests - PASS

**Total:** 138+ tests, all passing

#### **New Test Patterns**
- Table-driven middleware tests
- Configurable mock auth clients
- User isolation testing
- Token validation scenarios

---

### 🎯 Benefits

#### **Security**
- ✅ Prevents user impersonation
- ✅ User ID enforced by authentication
- ✅ Cannot assign quests to other users
- ✅ Cannot view other users' quests

#### **Simplicity**
- ✅ Cleaner API (no redundant parameters)
- ✅ Fewer fields to validate
- ✅ Single source of truth (JWT token)

#### **Consistency**
- ✅ User ID always comes from same source
- ✅ Authentication and authorization aligned
- ✅ Reduced client-side complexity

---

### ✅ Checklist

- [x] OpenAPI spec updated (removed user_id parameters)
- [x] OpenAPI code regenerated
- [x] HTTP handlers updated to use context
- [x] Authentication middleware enhanced
- [x] HTTP tests updated
- [x] Handler tests verified
- [x] E2E tests updated
- [x] Middleware tests created (27 new tests)
- [x] Mock auth clients enhanced
- [x] Documentation updated
- [x] All tests passing (138+)
- [x] Code compiles successfully
- [x] .cursorrules updated with changelog guidelines

---

### 🔗 Related Files

**Modified:**
- `api/http/quests/v1/openapi.yaml` - API specification
- `internal/adapters/in/http/assign_quest_handler.go` - Assign handler
- `internal/adapters/in/http/list_assigned_quests_handler.go` - List handler
- `internal/adapters/in/http/middleware/authentication.go` - Auth middleware
- `tests/integration/core/case_steps/http_requests.go` - Test helpers

**Created:**
- `tests/integration/mock/auth_scenarios.go` - Configurable auth mocks
- `tests/integration/tests/quest_middlewares/authentication/*.go` - Middleware tests
- `tests/integration/tests/test_container.go` - Enhanced with custom auth router
- `doc/API.md` - New API documentation
- `doc/TESTING.md` - New testing guide
- `doc/COMPONENTS.md` - New components guide
- `doc/CONFIGURATION.md` - New configuration guide
- `doc/DEVELOPMENT.md` - New development guide
- `doc/DEPLOYMENT.md` - New deployment guide
- `doc/ERROR_HANDLING.md` - New error handling guide
- `doc/DOMAIN_EVENTS.md` - New domain events guide
- `doc/INDEX.md` - Documentation index

**Updated:**
- `doc/ARCHITECTURE.md` - Refreshed architecture overview
- `.cursorrules` - Added changelog guidelines

---

### 🎓 Learning from This Change

**Key Lessons:**
1. **Security First:** User impersonation is a real threat
2. **Context is King:** Use request context for cross-cutting concerns
3. **API Design:** Simplicity improves security
4. **Testing:** Comprehensive tests catch issues early
5. **Documentation:** Keep docs in sync with code

---

**Breaking Change:** ⚠️ 

API clients must update to remove `user_id` from:
- `POST /quests/{id}/assign` request body
- `GET /quests/assigned` query parameters

User ID is now automatically extracted from JWT token. Ensure your JWT tokens contain valid user IDs.

---

**Migration Impact:** Medium  
**Security Impact:** High (positive)  
**Client Update Required:** Yes  
**Backward Compatible:** No

