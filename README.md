# Toucan

Toucan is a generic Learning Management System platform.

The project is planned as a Go backend with a React frontend, external identity
provider integration, and future support for mobile applications.

## Project Goals

- Provide a modular LMS foundation.
- Support external identity providers through OIDC/OAuth2.
- Work with existing user registries through identity mapping and future sync integrations.
- Support courses, enrollments, assessments, progress tracking, notifications, and reports.
- Keep the backend organized around business domains.
- Leave a clear path for future web, mobile, and enterprise integrations.

## Current Structure

```text
.
├── cmd/
│   └── toucan/
│       └── main.go
├── docs/
│   └── lms-system-plan.md
├── internal/
│   ├── assessments/
│   ├── content/
│   ├── courses/
│   ├── enrollments/
│   ├── identity/
│   ├── notifications/
│   ├── progress/
│   ├── reports/
│   ├── shared/
│   ├── tenants/
│   └── users/
├── go.mod
└── README.md
```

## Backend Organization

The Go backend uses a domain-oriented structure under `internal/`.

- `identity`: IDP integration, token validation, claims mapping, local identity resolution.
- `tenants`: tenant configuration, settings, feature flags, tenant-level isolation.
- `users`: LMS user profiles, user status, local user projection from external identities.
- `courses`: courses, versions, modules, lessons, publishing workflow.
- `content`: content items, attachments, file metadata, object storage integration.
- `enrollments`: enrollments, cohorts, groups, course assignments, learning paths.
- `assessments`: quizzes, questions, attempts, assignments, grading.
- `progress`: learner events, lesson progress, course progress, completions, certificates.
- `notifications`: email, in-app notifications, future mobile push notifications.
- `reports`: progress reports, completion reports, engagement reports, exports.
- `shared`: small cross-domain utilities that are truly generic.

Code should stay close to the domain it supports. Shared helpers should only move
to `shared` when more than one domain needs them and they have no domain-specific
behavior.

## Documentation

The system planning document is available at:

- `docs/lms-system-plan.md`

## Development

Run the backend entrypoint with:

```sh
go run ./cmd/toucan
```

## Storage Modes

The backend now supports two storage modes:

- `memory`: in-memory repositories for local development and tests
- `postgres`: PostgreSQL-backed repositories for the current production-ready path

Configuration is controlled through environment variables:

```sh
TOUCAN_HTTP_ADDR=:8080
TOUCAN_STORAGE_DRIVER=memory
```

To run against PostgreSQL:

```sh
TOUCAN_STORAGE_DRIVER=postgres
TOUCAN_POSTGRES_DSN=postgres://user:pass@localhost:5432/toucan?sslmode=disable
go run ./cmd/toucan
```

The PostgreSQL schema for the current `courses`, `sections`, and `content`
domains is created automatically on startup.
