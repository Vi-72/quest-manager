# Error Handling Guide - Quest Manager

## 🎯 Error Handling Strategy

Quest Manager uses structured error handling with **Problem Details (RFC 7807)** format for HTTP responses and domain-specific errors internally.

---

## 📊 Error Categories

### 1. Domain Errors (400 Bad Request)
**Source:** Business logic violations  
**Layer:** Domain  
**HTTP Status:** 400

**Examples:**
- Invalid quest status transition
- Quest already assigned
- Invalid difficulty value
- Reward out of range

**Handling:**
```go
// Domain layer
func (q *Quest) AssignTo(userID uuid.UUID) error {
    if q.status == StatusCompleted {
        return errors.New("cannot assign completed quest")
    }
    // ...
}

// Application layer wraps it
return errs.NewDomainValidationErrorWithCause("assignment", "failed to assign quest", err)

// HTTP response: 400 Bad Request
```

---

### 2. Not Found Errors (404 Not Found)
**Source:** Resource doesn't exist  
**Layer:** Application/Repository  
**HTTP Status:** 404

**Examples:**
- Quest ID not found
- Location ID not found

**Handling:**
```go
// Repository returns not found
quest, err := repo.GetByID(ctx, questID)
if err != nil {
    return errs.NewNotFoundErrorWithCause("quest", questID.String(), err)
}

// HTTP response: 404 Not Found
```

---

### 3. Validation Errors (400 Bad Request)
**Source:** Invalid input data  
**Layer:** HTTP/OpenAPI  
**HTTP Status:** 400

**Examples:**
- Missing required field
- Invalid UUID format
- Field value out of range
- Invalid enum value

**Handling:**
```go
// OpenAPI middleware validates automatically
// Returns 400 with field-specific error message
```

---

### 4. Authentication Errors (401 Unauthorized)
**Source:** Missing or invalid JWT token  
**Layer:** Middleware  
**HTTP Status:** 401

**Examples:**
- Missing Authorization header
- Expired JWT token
- Invalid token signature
- Malformed token

**Handling:**
```go
// Middleware
if errors.Is(err, auth.ErrTokenExpired) {
    problem := httperrors.NewProblem(
        http.StatusUnauthorized,
        "Token Expired",
        "JWT token has expired, please refresh your token",
    )
    problem.WriteResponse(w)
    return
}
```

---

### 5. Infrastructure Errors (500 Internal Server Error)
**Source:** External system failures  
**Layer:** Infrastructure  
**HTTP Status:** 500

**Examples:**
- Database connection failure
- Auth service unavailable
- Event publishing failure

**Handling:**
```go
if err := repo.Save(ctx, quest); err != nil {
    return errs.WrapInfrastructureError("failed to save quest", err)
}

// HTTP response: 500 Internal Server Error
// (Details hidden from client for security)
```

---

## 🔧 Error Types & Utilities

### Domain Error Types (`internal/pkg/errs/`)

#### `DomainValidationError`
```go
func NewDomainValidationError(domain, message string) error
func NewDomainValidationErrorWithCause(domain, message string, cause error) error

// Usage
return errs.NewDomainValidationError("assignment", "quest already assigned")
```

#### `NotFoundError`
```go
func NewNotFoundError(resourceType, resourceID string) error
func NewNotFoundErrorWithCause(resourceType, resourceID string, cause error) error

// Usage
return errs.NewNotFoundError("quest", questID.String())
```

#### `InfrastructureError`
```go
func WrapInfrastructureError(message string, cause error) error

// Usage
return errs.WrapInfrastructureError("failed to connect to database", err)
```

#### `ValueIsRequiredError`
```go
func NewValueIsRequiredError(fieldName string) error

// Usage
return errs.NewValueIsRequiredError("questRepository")
```

---

## 📝 HTTP Error Responses (RFC 7807)

### Problem Details Format

All HTTP errors use standardized format:

```json
{
  "type": "about:blank",
  "title": "Error Title",
  "status": 400,
  "detail": "Detailed error message"
}
```

### Error Response Examples

#### 400 Bad Request - Validation
```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 400,
  "detail": "Request validation failed for parameter 'difficulty': must be one of [easy, medium, hard]"
}
```

#### 400 Bad Request - Domain Logic
```json
{
  "type": "about:blank",
  "title": "Domain Validation Error",
  "status": 400,
  "detail": "assignment: failed to assign quest: quest cannot be assigned in status 'completed'"
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

#### 401 Unauthorized - Invalid Token
```json
{
  "type": "about:blank",
  "title": "Authentication Failed",
  "status": 401,
  "detail": "Invalid or malformed authentication token"
}
```

#### 404 Not Found
```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "quest with id '550e8400-e29b-41d4-a716-446655440000' not found"
}
```

#### 409 Conflict
```json
{
  "type": "about:blank",
  "title": "Conflict",
  "status": 409,
  "detail": "Resource already exists with this identifier"
}
```

#### 500 Internal Server Error
```json
{
  "type": "about:blank",
  "title": "Internal Server Error",
  "status": 500,
  "detail": "An unexpected error occurred. Please try again later."
}
```

**Security Note:** 500 errors don't expose internal details to clients.

---

## 🔄 Error Flow

### Domain Error → HTTP Response

```
Domain Layer
    ↓ errors.New("business rule violation")
Application Layer
    ↓ errs.NewDomainValidationError(...)
HTTP Handler
    ↓ return nil, err
Error Middleware
    ↓ Convert to Problem Details
HTTP Response
    ↓ 400 Bad Request + JSON
Client
```

---

## 🛡️ Error Handling Best Practices

### 1. Always Wrap Errors
```go
// Good: add context
if err := repo.Save(ctx, quest); err != nil {
    return errs.WrapInfrastructureError("failed to save quest", err)
}

// Bad: return raw error
return err
```

### 2. Use Specific Error Types
```go
// Good: specific error type
return errs.NewNotFoundError("quest", id)

// Bad: generic error
return errors.New("not found")
```

### 3. Log Errors Appropriately
```go
// Good: log with context
log.Printf("ERROR failed to save quest quest_id=%s error=%v", questID, err)

// Bad: silent failure
return err
```

### 4. Don't Expose Internal Details
```go
// Good: generic message for infrastructure errors
return httperrors.NewInternalServerError()

// Bad: expose stack trace
return errors.New(fmt.Sprintf("SQL error: %v", sqlErr))
```

### 5. Validate Early
```go
// Good: validate in domain constructor
func NewQuest(...) (*Quest, error) {
    if reward < 1 || reward > 5 {
        return nil, errors.New("reward must be 1-5")
    }
}

// Bad: validate in handler
if quest.Reward < 1 {
    return error
}
```

---

## 🧪 Testing Error Scenarios

### Test Domain Errors
```go
func TestQuest_AssignTo_InvalidStatus(t *testing.T) {
    quest := createQuestWithStatus(StatusCompleted)
    
    err := quest.AssignTo(userID)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot assign")
}
```

### Test HTTP Errors
```go
func (s *Suite) TestAssignQuestHTTP_QuestNotFound() {
    req := AssignQuestHTTPRequest(nonExistentID)
    resp := ExecuteHTTPRequest(ctx, router, req)
    
    assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    assert.Contains(t, resp.Body, "not found")
}
```

### Test Authentication Errors
```go
func (s *Suite) TestCreateQuestWithExpiredToken() {
    router := NewRouterWithExpiredTokenAuth()
    
    resp := ExecuteHTTPRequest(ctx, router, createReq)
    
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
    assert.Contains(t, resp.Body, "expired")
}
```

---

## 📋 Common Error Scenarios

### Scenario 1: Quest Not Found
```
User Action: GET /quests/invalid-uuid
    ↓
Handler: GetQuestByID
    ↓
Repository: GetByID returns not found
    ↓
Handler: errs.NewNotFoundError("quest", id)
    ↓
Error Middleware: Convert to Problem Details
    ↓
Response: 404 Not Found
```

### Scenario 2: Invalid Status Transition
```
User Action: PATCH /quests/{id}/status {"status": "completed"}
    ↓
Handler: ChangeQuestStatus
    ↓
Domain: quest.ChangeStatus() validates transition
    ↓
Domain: returns "invalid transition from 'created' to 'completed'"
    ↓
Handler: errs.NewDomainValidationError(...)
    ↓
Response: 400 Bad Request
```

### Scenario 3: Database Connection Lost
```
User Action: POST /quests
    ↓
Handler: CreateQuest
    ↓
Repository: Save() fails (db connection error)
    ↓
Handler: errs.WrapInfrastructureError("failed to save", err)
    ↓
Response: 500 Internal Server Error
    ↓
Logged: ERROR with full stack trace
    ↓
Client sees: Generic error message (security)
```

---

## 🔗 Related

- [API Documentation](API.md) - API error responses
- [Testing](TESTING.md) - Testing error scenarios
- [Components](COMPONENTS.md) - Error handling in components

---

**Principle:** Errors should be informative for developers but safe for users.

