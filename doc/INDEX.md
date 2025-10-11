# Quest Manager - Documentation Index

## 📚 Complete Documentation Guide

Welcome to Quest Manager documentation! This index will help you find the information you need.

---

## 🎯 Quick Navigation

### For API Consumers
- [**API Documentation**](API.md) - REST API endpoints, authentication, examples
- [**Error Handling**](ERROR_HANDLING.md) - Error responses and troubleshooting

### For Developers
- [**Development Guide**](DEVELOPMENT.md) - Setup, workflows, best practices
- [**Testing Guide**](TESTING.md) - Test strategies, patterns, running tests
- [**Components**](COMPONENTS.md) - System components and architecture
- [**Domain Events**](DOMAIN_EVENTS.md) - Event system and usage

### For DevOps
- [**Configuration**](CONFIGURATION.md) - Environment variables, settings
- [**Deployment**](DEPLOYMENT.md) - Docker, Kubernetes, production setup
- [**Architecture**](ARCHITECTURE.md) - High-level system design

### For Project Management
- [**Changelog**](changelog/) - Version history and changes
- [**README**](../README.md) - Project overview

---

## 📖 Documentation by Topic

### Architecture & Design

| Document                          | Description                            | Audience               |
|-----------------------------------|----------------------------------------|------------------------|
| [Architecture](ARCHITECTURE.md)   | System design, DDD, Clean Architecture | Developers, Architects |
| [Components](COMPONENTS.md)       | Component breakdown, dependencies      | Developers             |
| [Domain Events](DOMAIN_EVENTS.md) | Event system, event catalog            | Developers             |

### API & Integration

| Document                                 | Description                    | Audience                     |
|------------------------------------------|--------------------------------|------------------------------|
| [API Documentation](API.md)              | REST endpoints, authentication | API Consumers, Frontend Devs |
| [Error Handling](ERROR_HANDLING.md)      | Error types, Problem Details   | API Consumers, Developers    |
| OpenAPI Spec                             | Machine-readable API spec      | Tools, Generators            |

### Development

| Document                            | Description                       | Audience           |
|-------------------------------------|-----------------------------------|--------------------|
| [Development Guide](DEVELOPMENT.md) | Local setup, workflows, debugging | Developers         |
| [Testing Guide](TESTING.md)         | Test pyramid, patterns, coverage  | Developers, QA     |
| [Configuration](CONFIGURATION.md)   | Environment variables, settings   | Developers, DevOps |

### Operations

| Document                          | Description                   | Audience    |
|-----------------------------------|-------------------------------|-------------|
| [Deployment](DEPLOYMENT.md)       | Docker, K8s, production setup | DevOps, SRE |
| [Configuration](CONFIGURATION.md) | Env-specific configs          | DevOps      |

### Version History

| Document                                                          | Description                | Audience |
|-------------------------------------------------------------------|----------------------------|----------|
| [Changelog - Auth](changelog/1_CHANGELOG_AUTH.md)                 | JWT Authentication feature | All      |
| [Changelog - User Context](changelog/2_CHANGELOG_USER_CONTEXT.md) | User from token feature    | All      |

---

## 🎓 Learning Paths

### Path 1: New Developer Onboarding

1. **Start here:** [README](../README.md) - Project overview
2. **Understand design:** [Architecture](ARCHITECTURE.md) - System design
3. **Setup environment:** [Development](DEVELOPMENT.md) - Local setup
4. **Explore API:** [API Documentation](API.md) - What the system does
5. **Understand components:** [Components](COMPONENTS.md) - How it works
6. **Learn testing:** [Testing](TESTING.md) - How to test
7. **Start coding:** Pick a task and follow development guide

**Time estimate:** 4-6 hours

---

### Path 2: API Consumer Integration

1. **API basics:** [API Documentation](API.md) - Endpoints and auth
2. **Authentication:** Section in API doc - Get JWT token
3. **Error handling:** [Error Handling](ERROR_HANDLING.md) - Handle errors
4. **Try examples:** Use cURL examples from API doc
5. **Implement client:** Use OpenAPI spec for code generation

**Time estimate:** 1-2 hours

---

### Path 3: DevOps Deployment

1. **Configuration:** [Configuration](CONFIGURATION.md) - Environment setup
2. **Deployment:** [Deployment](DEPLOYMENT.md) - Deploy options
3. **Architecture:** [Architecture](ARCHITECTURE.md) - System dependencies
4. **Monitoring:** Check deployment guide for health checks
5. **Troubleshooting:** [Error Handling](ERROR_HANDLING.md) - Debug issues

**Time estimate:** 2-3 hours

---

## 🔍 Finding Information

### By Task

**"I want to..."**

- **...understand the system:** → [Architecture](ARCHITECTURE.md)
- **...use the API:** → [API Documentation](API.md)
- **...develop a feature:** → [Development](DEVELOPMENT.md)
- **...write tests:** → [Testing](TESTING.md)
- **...deploy the app:** → [Deployment](DEPLOYMENT.md)
- **...configure settings:** → [Configuration](CONFIGURATION.md)
- **...handle errors:** → [Error Handling](ERROR_HANDLING.md)
- **...understand events:** → [Domain Events](DOMAIN_EVENTS.md)
- **...see what changed:** → [Changelog](changelog/)

### By Component

| Component | Documentation |
|-----------|---------------|
| HTTP Handlers | [Components - HTTP Adapters](COMPONENTS.md#4-http-adapters-internaladaptersinhttp) |
| Domain Models | [Components - Domain Layer](COMPONENTS.md#1-domain-layer-internalcoredomain) |
| Use Cases | [Components - Application Layer](COMPONENTS.md#2-application-layer-internalcoreapplication) |
| Repositories | [Components - Infrastructure](COMPONENTS.md#5-infrastructure-adapters-internaladaptersout) |
| Middleware | [Components - Middleware](COMPONENTS.md#middleware-middleware) |
| Auth System | [API - Authentication](API.md#-authentication) |
| Event System | [Domain Events](DOMAIN_EVENTS.md) |

---

## 🆘 Troubleshooting Guide

| Problem | Documentation | Section |
|---------|---------------|---------|
| Tests failing | [Testing](TESTING.md) | Debugging Tests |
| API errors | [Error Handling](ERROR_HANDLING.md) | Error Categories |
| Configuration issues | [Configuration](CONFIGURATION.md) | Common Issues |
| Deployment problems | [Deployment](DEPLOYMENT.md) | Rollback Procedure |
| Auth not working | [API](API.md) | Authentication |
| Database errors | [Development](DEVELOPMENT.md) | Common Issues |

---

## 📝 Documentation Maintenance

### When to Update

**Update immediately:**
- New API endpoints added
- Breaking changes to API
- New configuration options
- Architecture changes

**Update regularly:**
- New features
- Bug fixes (if significant)
- Performance improvements
- Security updates

### Documentation Standards

1. **Use Markdown:** All docs in `.md` format
2. **Clear structure:** Use headers, lists, code blocks
3. **Examples:** Provide concrete examples
4. **Keep updated:** Update with code changes
5. **Version:** Note last updated date

---

## 🔗 External Resources

### Related Services
- [Quest Auth Service](https://github.com/Vi-72/quest-auth) - Authentication service
- OpenAPI Specification - `api/http/quests/v1/openapi.yaml`

### Technologies
- [Go](https://go.dev/doc/) - Programming language
- [GORM](https://gorm.io/docs/) - ORM for database
- [PostgreSQL](https://www.postgresql.org/docs/) - Database
- [gRPC](https://grpc.io/docs/) - RPC framework

### Architecture Patterns
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [DDD](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [CQRS](https://martinfowler.com/bliki/CQRS.html)

---

## 🗺️ Documentation Map

```
doc/
├── INDEX.md (you are here)          # Documentation index
├── ARCHITECTURE.md                   # System design
├── API.md                           # API reference
├── COMPONENTS.md                    # Component details
├── CONFIGURATION.md                 # Settings guide
├── DEVELOPMENT.md                   # Dev guide
├── DEPLOYMENT.md                    # Deployment guide
├── TESTING.md                       # Testing guide
├── ERROR_HANDLING.md                # Error guide
├── DOMAIN_EVENTS.md                 # Events guide
└── changelog/                       # Version history
    ├── 1_CHANGELOG_AUTH.md         # Auth feature
    └── 2_CHANGELOG_USER_CONTEXT.md # User context feature
```

---

## 💡 Quick Tips

### For New Developers
1. Start with [Architecture](ARCHITECTURE.md) to understand system design
2. Read [Components](COMPONENTS.md) to learn component structure
3. Follow [Development](DEVELOPMENT.md) to setup local environment
4. Check [Testing](TESTING.md) before writing code

### For API Users
1. [API Documentation](API.md) has everything you need
2. Test with cURL examples first
3. Handle errors properly (see [Error Handling](ERROR_HANDLING.md))

### For DevOps
1. [Deployment](DEPLOYMENT.md) covers all deployment scenarios
2. [Configuration](CONFIGURATION.md) lists all settings
3. Check [Architecture](ARCHITECTURE.md) for system dependencies

---

## 📮 Feedback

Found an error in documentation? Have suggestions?
- Create an issue
- Submit a PR with improvements
- Contact the team

---

## 📊 Documentation Stats

- **Total Documents:** 10
- **Total Pages:** ~50+ (estimated)
- **Code Examples:** 100+
- **Diagrams:** 5+
- **Last Updated:** October 9, 2025

---

**Welcome to Quest Manager!** Start exploring the documentation using the links above. 🚀

