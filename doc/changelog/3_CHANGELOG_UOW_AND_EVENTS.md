# UnitOfWork Read Consistency and Event Publisher Fallback - Changelog

## Version 1.5.1 - UoW consistency and event publishing improvements

### ✨ New Features
- Read-your-writes consistency for repository reads inside active transactions
- Event publisher synchronous publishing now works even without UoW in context (fallback)

### 🔧 Technical Changes

#### New/Updated Components
1. **Quest Repository** (`internal/adapters/out/postgres/questrepo/repository.go`)
   - Added `dbWithCtx(ctx)` to route reads via `Tx()` when `InTx()==true`
   - Updated read methods to use transactional context where applicable

2. **Location Repository** (`internal/adapters/out/postgres/locationrepo/repository.go`)
   - Added `dbWithCtx(ctx)` and switched reads to transactional when `InTx()==true`

3. **Event Repository (Publisher)** (`internal/adapters/out/postgres/eventrepo/repository.go`)
   - `Publish(ctx, events...)` now prefers UoW from context, but falls back to creating a new UoW when none present

4. **Command Executor** (`internal/core/application/usecases/commands/command_executor.go`)
   - Fixed commit context: `Commit(ctxWithUow)` to ensure consistent UoW propagation

5. **Tests DI (integration)** (`tests/integration/tests/test_container.go`)
   - Query handlers now receive `UnitOfWorkFactory`

6. **Mocks DI (contracts)** (`tests/contracts/mocks/contract_di_container.go`)
   - Handlers share the same in-memory repos via mock UoW for consistent state

### 📝 Configuration
- No configuration changes

### 🔒 Security
- No changes

### 📊 Error Responses
- No changes to API error contract

### 🔄 Migration Guide
- No DB migrations required
- Consumers should see consistent reads inside transactional flows

### 📚 Dependencies
- No dependency changes

### 🧪 Testing Notes
- Contracts: all passing
- Repository integration tests: passing
- HTTP/E2E integration tests require running Postgres and mock/real auth

### 🎯 Future Improvements
- Add negative-path tests for query handlers (UoW factory errors, repo errors)
- Improve coverage for `ports/uow_context.go` and `pkg/errs`

### ✅ Checklist
- [x] Read-your-writes for repos inside transactions
- [x] Event publisher fallback when no UoW in context
- [x] Fixed commit context in `CommandExecutor`
- [x] Updated integration and contract DI

### 🔗 Related
- `internal/adapters/out/postgres/questrepo/repository.go`
- `internal/adapters/out/postgres/locationrepo/repository.go`
- `internal/adapters/out/postgres/eventrepo/repository.go`
- `internal/core/application/usecases/commands/command_executor.go`
- `tests/integration/tests/test_container.go`
- `tests/contracts/mocks/contract_di_container.go`

---
**Breaking Change:** None
