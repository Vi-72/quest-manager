# Quest Manager - Project Brief

## 🎯 Project Overview

**Quest Manager** is a backend HTTP API service for creating and managing location-based quests with advanced geospatial capabilities.

## 🎪 Core Purpose

Quest Manager serves as the central hub for a quest-based application ecosystem, enabling:

- **Quest Lifecycle Management**: Complete CRUD operations for quests
- **Geospatial Operations**: Location-based quest discovery and assignment
- **User Context Integration**: Secure user identification and authorization
- **Event-Driven Architecture**: Real-time tracking of quest state changes

## 🏗️ Architectural Foundation

### Design Philosophy
- **Clean Architecture**: Clear separation of concerns across layers
- **Domain-Driven Design**: Business logic centered around quest and location domains
- **Event-Driven**: State changes communicated through domain events
- **Security-First**: JWT authentication with external auth service integration

### Core Principles
1. **Domain-Centric**: Business logic isolated in domain layer
2. **Dependency Inversion**: Outer layers depend on inner layers
3. **CQRS**: Separation of command and query operations
4. **Hexagonal Architecture**: Ports and adapters for external integrations

## 🎯 Key Requirements

### Functional Requirements
- ✅ **Quest Management**: Create, read, update, assign quests
- ✅ **Location Services**: Geospatial search within radius
- ✅ **User Authentication**: JWT-based security with external auth service
- ✅ **Status Workflow**: Quest lifecycle management (created → posted → assigned → completed)
- ✅ **Event Tracking**: Domain events for audit and integration

### Non-Functional Requirements
- **Performance**: <100ms response time for typical operations
- **Scalability**: Horizontal scaling capability
- **Security**: JWT authentication, SQL injection prevention
- **Reliability**: Transactional consistency, graceful error handling
- **Testability**: >75% code coverage, comprehensive test suite

## 🎪 Domain Model

### Core Aggregates
1. **Quest**: Central business entity with lifecycle management
2. **Location**: Geospatial entity with coordinate-based operations
3. **GeoCoordinate**: Value object for precise location calculations

### Key Business Rules
- Quests have mandatory target and execution locations
- Quest status transitions follow defined workflow
- Location coordinates must be valid (latitude: -90 to 90, longitude: -180 to 180)
- User context extracted from JWT tokens, not request parameters

## 🔧 Technical Constraints

### Technology Stack
- **Language**: Go 1.23+
- **Database**: PostgreSQL with spatial extensions
- **Authentication**: External gRPC auth service
- **API**: REST with OpenAPI 3.0 specification
- **Architecture**: Clean Architecture + DDD

### Integration Requirements
- **Quest Auth Service**: gRPC client for JWT validation
- **PostgreSQL**: Primary data store with spatial indexing
- **OpenAPI**: Code generation from specification

## 🎯 Success Criteria

### Quality Metrics
- **Code Coverage**: >75% (Current: 75.6%)
- **Test Suite**: 110+ tests across all layers
- **Performance**: <100ms average response time
- **Reliability**: Zero data loss, transactional consistency

### Business Value
- **Developer Experience**: Clean, testable, maintainable codebase
- **System Integration**: Event-driven architecture for extensibility
- **Security**: Production-ready authentication and authorization
- **Scalability**: Foundation for future microservices evolution

## 🚀 Project Scope

### In Scope
- Quest CRUD operations
- Location-based quest discovery
- JWT authentication integration
- Domain event system
- Comprehensive test coverage
- OpenAPI documentation

### Out of Scope (Current Version)
- User management (handled by external auth service)
- Real-time notifications
- File uploads
- Advanced analytics
- Multi-tenancy

## 📊 Current Status

**Version**: 1.5.0  
**Status**: Production Ready ✅  
**Last Updated**: October 9, 2025

### Achievements
- ✅ Complete Clean Architecture implementation
- ✅ Domain-driven design with rich domain models
- ✅ CQRS pattern with command/query separation
- ✅ Event-driven architecture with domain events
- ✅ Comprehensive test suite (110+ tests)
- ✅ JWT authentication with external service
- ✅ Geospatial operations with optimized queries
- ✅ OpenAPI specification with code generation

### Quality Metrics
- **Test Coverage**: 75.6%
- **Architecture Compliance**: Clean Architecture + DDD
- **Security**: JWT authentication, input validation
- **Performance**: Optimized database queries with spatial indexing
- **Documentation**: Comprehensive API and architecture docs

---

**This document serves as the foundation for all other Memory Bank files and defines the core scope and requirements of the Quest Manager project.**

