# Domain Events - Quest Manager

## 🎯 Overview

Quest Manager uses **Domain Events** to communicate changes within the system following DDD principles.

---

## 📋 Event Catalog

### Quest Events

#### `quest.created`
**Trigger:** New quest is created  
**Data:**
```json
{
  "quest_id": "uuid",
  "creator": "user-id",
  "title": "string",
  "difficulty": "easy|medium|hard",
  "status": "created"
}
```

---

#### `quest.assigned`
**Trigger:** Quest is assigned to a user  
**Data:**
```json
{
  "quest_id": "uuid",
  "assignee": "user-id",
  "previous_status": "created|posted",
  "new_status": "assigned"
}
```

---

#### `quest.status_changed`
**Trigger:** Quest status changes  
**Data:**
```json
{
  "quest_id": "uuid",
  "from_status": "previous-status",
  "to_status": "new-status",
  "changed_by": "user-id"
}
```

---

### Location Events

#### `location.created`
**Trigger:** New location is created  
**Data:**
```json
{
  "location_id": "uuid",
  "latitude": 55.7558,
  "longitude": 37.6173,
  "address": "optional-address"
}
```

---

#### `location.updated`
**Trigger:** Location coordinates or address updated  
**Data:**
```json
{
  "location_id": "uuid",
  "previous_coordinates": {...},
  "new_coordinates": {...}
}
```

---

## 🔄 Event Flow

### Event Lifecycle

```
1. Domain Operation
   ↓
2. Domain Model adds event
   quest.AddDomainEvent(NewQuestCreatedEvent(...))
   ↓
3. Handler publishes events (in transaction)
   eventPublisher.Publish(ctx, quest.GetDomainEvents()...)
   ↓
4. Event persisted to database
   INSERT INTO events (...)
   ↓
5. Transaction commits
   ↓
6. Events cleared from aggregate
   quest.ClearDomainEvents()
```

### Publishing Strategy

**Transactional Publishing:**
- Events published within same transaction as aggregate
- If transaction rolls back, events are not persisted
- Ensures consistency between aggregate state and events

**Async Processing:**
- Events persisted synchronously (in transaction)
- Future: Async consumers can process events from event store

---

## 🏗️ Event Architecture

### Domain Event Interface
```go
type DomainEvent interface {
    EventID() uuid.UUID
    OccurredAt() time.Time
    AggregateID() uuid.UUID
    EventName() string
}
```

### Base Event
```go
type BaseEvent struct {
    eventID     uuid.UUID
    occurredAt  time.Time
    aggregateID uuid.UUID
}
```

### Specific Events

**Quest Created:**
```go
type QuestCreatedEvent struct {
    BaseEvent
    creator    string
    title      string
    difficulty string
    status     string
}
```

**Quest Assigned:**
```go
type QuestAssignedEvent struct {
    BaseEvent
    assignee       uuid.UUID
    previousStatus Status
    newStatus      Status
}
```

---

## 📦 Event Publisher

### Interface
```go
type EventPublisher interface {
    Publish(ctx context.Context, events ...ddd.DomainEvent) error
    PublishAsync(events ...ddd.DomainEvent)
}
```

### Implementation (`eventrepo/repository.go`)

**Features:**
- Goroutine pool for async publishing
- Transactional support
- Event persistence to PostgreSQL
- JSON serialization of event data

**Configuration:**
```go
// Create publisher with goroutine limit
publisher := eventrepo.NewRepository(tracker, 10)
```

---

## 🗄️ Event Storage

### Database Schema
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    event_name VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    event_data JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_events_aggregate ON events(aggregate_id);
CREATE INDEX idx_events_name ON events(event_name);
CREATE INDEX idx_events_occurred_at ON events(occurred_at);
```

### Event Record
```json
{
  "id": "event-uuid",
  "event_name": "quest.created",
  "aggregate_id": "quest-uuid",
  "aggregate_type": "Quest",
  "event_data": {
    "creator": "user-id",
    "title": "Find the treasure",
    "difficulty": "medium",
    "status": "created"
  },
  "occurred_at": "2025-10-09T10:30:00Z",
  "created_at": "2025-10-09T10:30:01Z"
}
```

---

## 🎯 Event Usage Patterns

### Pattern 1: Single Event
```go
// Create quest
quest := quest.NewQuest(...)

// Event is added automatically
events := quest.GetDomainEvents()
// Contains: [quest.created]

// Publish
eventPublisher.Publish(ctx, events...)

// Clear after commit
quest.ClearDomainEvents()
```

### Pattern 2: Multiple Events
```go
// Assign quest
quest.AssignTo(userID)

// Two events added:
// 1. quest.assigned
// 2. quest.status_changed (created → assigned)

events := quest.GetDomainEvents()
// Contains: [quest.assigned, quest.status_changed]
```

### Pattern 3: Conditional Events
```go
// Change status
if err := quest.ChangeStatus(newStatus); err != nil {
    return err
}

// Event only added if status actually changed
if oldStatus != newStatus {
    quest.AddDomainEvent(NewStatusChangedEvent(...))
}
```

---

## 🧪 Testing Events

### Test Event Generation
```go
func TestQuest_NewQuest_DomainEvents(t *testing.T) {
    quest := quest.NewQuest(...)
    
    events := quest.GetDomainEvents()
    
    assert.Len(t, events, 1)
    assert.Equal(t, "quest.created", events[0].EventName())
}
```

### Test Event Publishing
```go
func (s *Suite) TestAssignQuest_PublishesEvents() {
    // Act
    handler.Handle(ctx, assignCmd)
    
    // Assert: event persisted in database
    events := eventStorage.GetEventsByAggregate(ctx, questID)
    assert.Contains(t, events, "quest.assigned")
}
```

### Test Event Data
```go
func (s *Suite) TestQuestCreatedEvent_Data() {
    // Retrieve event from storage
    event := eventStorage.GetEvent(ctx, eventID)
    
    // Assert event data
    assert.Equal(t, "quest.created", event.EventName)
    assert.Equal(t, questID, event.AggregateID)
    
    // Parse event data
    var data QuestCreatedEventData
    json.Unmarshal(event.EventData, &data)
    assert.Equal(t, "medium", data.Difficulty)
}
```

---

## 🔮 Future Event Consumers

### Event Sourcing (Potential)
- Rebuild aggregate state from events
- Event replay for debugging
- Audit trail

### Event Notifications (Potential)
- WebSocket notifications to clients
- Email notifications on quest assignment
- Push notifications for mobile apps

### Analytics (Potential)
- Track quest creation patterns
- Analyze assignment rates
- Monitor status transitions

### Integration (Potential)
- Publish events to message broker (RabbitMQ, Kafka)
- Trigger external workflows
- Sync with other microservices

---

## 📊 Event Metrics

### Current Events
- `quest.created` - ~100% of quest creations
- `quest.assigned` - ~80% of quests
- `quest.status_changed` - ~5-10 per quest lifecycle
- `location.created` - 2x per quest (target + execution)

### Event Volume (estimated)
- **Low traffic:** ~10 events/minute
- **Medium traffic:** ~100 events/minute
- **High traffic:** ~1000 events/minute

### Performance
- Event publishing: <5ms (sync)
- Event persistence: <10ms (with transaction)
- Goroutine pool prevents overload

---

## ⚙️ Event Configuration

### Goroutine Limit
```bash
# Control concurrent event processing
export EVENT_GOROUTINE_LIMIT=20

# Default: 10
```

**Tuning:**
- Low traffic: 5-10 goroutines
- Medium traffic: 10-20 goroutines
- High traffic: 20-50 goroutines

### Event Retention (Future)
```sql
-- Archive old events
CREATE TABLE events_archive AS 
SELECT * FROM events WHERE occurred_at < NOW() - INTERVAL '90 days';

DELETE FROM events WHERE occurred_at < NOW() - INTERVAL '90 days';
```

---

## 🔧 Troubleshooting

### Events Not Persisting
**Possible causes:**
1. Transaction rollback
2. Event publisher not injected
3. Database connection issue

**Debug:**
```go
log.Printf("Publishing %d events", len(events))
err := publisher.Publish(ctx, events...)
if err != nil {
    log.Printf("ERROR publishing events: %v", err)
}
```

### Duplicate Events
**Cause:** Events not cleared after publishing

**Solution:**
```go
// Always clear after successful commit
if err := unitOfWork.Commit(ctx); err == nil {
    aggregate.ClearDomainEvents()
}
```

### Event Data Corruption
**Cause:** JSON serialization issues

**Solution:**
- Ensure event data is JSON-serializable
- Use primitive types in event data
- Test event serialization/deserialization

---

## 🔗 Related

- [Components](COMPONENTS.md) - Event publisher component
- [Architecture](ARCHITECTURE.md) - Event-driven architecture
- [Testing](TESTING.md) - Testing events

---

**Principle:** Events represent facts about what happened in the system. They are immutable and append-only.

