# Generic LMS System Plan

## Overview

This document outlines a plan for building a generic Learning Management System
with a Go backend, React frontend, external identity provider integration, and a
future path for mobile applications.

The system should delegate authentication to an external identity provider
while keeping LMS-specific authorization, course data, progress tracking, and
business rules inside the LMS.

The core design principle is:

> The LMS should not own user authentication. It should consume trusted identity
> claims from an IDP, map those identities to local LMS users, and enforce
> application permissions internally.

## Target Architecture

```text
Web App / Future Mobile Apps
        |
        | OIDC Authorization Code + PKCE
        v
Identity Provider
        |
        | JWT / ID token / access token
        v
Go Backend API
        |
        +-- LMS Core Modules
        +-- Authorization / RBAC
        +-- Course & Content Management
        +-- Enrollment
        +-- Assessments
        +-- Progress Tracking
        +-- Notifications
        +-- Reporting
        |
        v
PostgreSQL / Object Storage / Redis / Queue
```

## Recommended Stack

### Backend

- Go
- REST API to start; GraphQL can be considered later if client data needs become complex
- PostgreSQL as the primary database
- Redis for caching, rate limits, and background job coordination
- S3-compatible object storage for files, videos, thumbnails, and documents
- Background workers in Go for emails, reports, imports, exports, and notification jobs
- OpenAPI for API documentation and frontend/mobile client generation

### Frontend

- React with TypeScript
- React Router or TanStack Router
- TanStack Query for server state
- Component library such as MUI, Mantine, shadcn/ui, or a local design system
- OIDC client library for login, token handling, and callback processing
- Permission-aware route guards and UI controls

### Identity Provider

The IDP should support OpenID Connect and OAuth2.

Suitable options include:

- Keycloak
- Auth0
- Okta
- Microsoft Entra ID
- AWS Cognito
- FusionAuth

Use Authorization Code Flow with PKCE for both the web frontend and future
mobile applications.

## Major System Modules

## 1. Tenant And Organization Management

If the LMS may serve multiple institutions, companies, schools, or departments,
design for tenancy from the beginning.

Core entities:

- Tenant
- Organization unit
- Department or group
- Branding settings
- IDP connection settings
- Feature flags

Even if the first version is deployed for a single organization, adding
`tenant_id` to core tables early will avoid painful migrations later.

## 2. Identity Mapping

The LMS should maintain a local user record as an application-level projection
of the external identity.

Example user fields:

```text
users
- id
- tenant_id
- external_subject
- idp_issuer
- email
- display_name
- status
- last_login_at
- created_at
- updated_at
```

Use `idp_issuer + external_subject` as the stable identity key. Do not rely only
on email, because email addresses can change or collide across identity sources.

## 3. Authorization

Authentication should come from the IDP. Authorization should mostly live inside
the LMS.

Recommended authorization layers:

- Global roles: platform admin, support admin
- Tenant roles: tenant admin, tenant manager
- Organization roles: instructor, learner, reviewer, manager
- Course roles: course owner, instructor, teaching assistant, enrolled learner
- Permission checks: `course.read`, `course.manage`, `assessment.grade`, `report.view`

Avoid scattering raw role checks through the codebase. Prefer central permission
checks such as:

```go
Can(user, "course.manage", courseID)
Can(user, "assessment.grade", assessmentID)
Can(user, "report.view", tenantID)
```

## 4. Course Management

Core entities:

```text
courses
course_versions
modules
lessons
content_items
attachments
course_categories
course_tags
```

Course versioning is worth considering early. A learner may need to complete the
version they enrolled in, while new learners should receive the latest published
version.

## 5. Enrollment

Core entities:

```text
enrollments
cohorts
groups
learning_paths
course_assignments
```

Enrollment sources may include:

- Manual admin enrollment
- Self-enrollment
- Bulk CSV import
- API import from an existing registry
- IDP group-based enrollment
- HRIS or SIS integration later

## 6. Learning Progress

Track progress separately from content.

Core entities:

```text
lesson_progress
course_progress
activity_events
completion_records
certificates
```

Use append-only events for important learner actions:

```text
learner_started_lesson
learner_completed_lesson
learner_submitted_assessment
learner_passed_course
learner_downloaded_certificate
```

Maintain summary tables for fast reads in dashboards and reports.

## 7. Assessments

Initial entities:

```text
quizzes
questions
question_options
quiz_attempts
answers
grades
rubrics
assignments
submissions
```

Initial question types:

- Multiple choice
- Multiple select
- True/false
- Short answer
- File upload
- Essay

Future assessment capabilities:

- Randomized question banks
- Timed assessments
- Proctoring integration
- Plagiarism detection
- SCORM or xAPI support

## 8. Content Delivery

Supported content types:

- Rich text
- Video
- PDF
- File download
- External link
- Embedded content
- Quiz
- Assignment
- Live session link

Store content metadata in PostgreSQL and binary files in object storage.

For video, avoid building a custom streaming platform unless required. Consider
using:

- Vimeo Enterprise
- Cloudflare Stream
- AWS MediaConvert with CloudFront
- Mux

## 9. Notifications

Notification channels:

- Email
- In-app notifications
- Future mobile push notifications

Notification events:

- Course assigned
- Deadline approaching
- Assessment graded
- Certificate issued
- Announcement posted
- Instructor feedback received

Use a background job queue so notifications do not block API requests.

## 10. Reporting

Initial reports:

- Learner progress
- Course completion
- Assessment scores
- Enrollment counts
- Active and inactive users
- Course engagement
- Certificates issued

For heavier analytics later, stream events into a warehouse or analytics
database.

## Backend Service Design

Start with a modular monolith in Go rather than microservices.

Suggested structure:

```text
backend/
  cmd/
    api/
    worker/
  internal/
    auth/
    tenants/
    users/
    courses/
    enrollments/
    content/
    assessments/
    progress/
    notifications/
    reports/
    storage/
    database/
  migrations/
  api/
    openapi.yaml
```

Reasons to start with a modular monolith:

- Easier transactions
- Simpler deployment
- Faster product iteration
- Less operational overhead
- Easier debugging
- Services can still be extracted later if module boundaries are kept clean

## API Design

Use versioned APIs.

Example endpoints:

```text
/api/v1/me
/api/v1/courses
/api/v1/courses/{courseId}
/api/v1/courses/{courseId}/modules
/api/v1/enrollments
/api/v1/assessments
/api/v1/progress
/api/v1/reports
```

API principles:

- Backend enforces all permissions
- Frontend permissions are only for user experience
- Use pagination on all list endpoints
- Use stable IDs, preferably UUID or ULID
- Resolve `tenant_id` from backend context, not trusted frontend input
- Generate OpenAPI clients for frontend and mobile apps when possible

## Authentication Flow

For the React web app:

1. User opens the LMS.
2. React redirects the user to the IDP using OIDC Authorization Code Flow with PKCE.
3. User authenticates with the IDP.
4. React receives tokens through the OIDC callback.
5. React calls the Go backend with the access token.
6. Backend validates the JWT signature using the IDP JWKS endpoint.
7. Backend validates issuer, audience, expiry, and required claims.
8. Backend resolves the local LMS user by `issuer + subject`.
9. Backend applies tenant, role, and permission checks.
10. Backend returns LMS data.

For future mobile apps, use the same OIDC Authorization Code Flow with PKCE and
a native redirect URI.

## IDP Integration Strategy

Support these integration modes over time.

### Single IDP Per Tenant

This is the simplest model. Each tenant configures one OIDC provider.

### Multiple IDPs Per Tenant

Useful when one organization has multiple user registries.

### IDP Group Claim Mapping

Map external groups to LMS roles, cohorts, or permissions.

Example:

```text
IDP group: engineering-learners
LMS role: learner
LMS cohort: Engineering 2026
```

### Just-In-Time User Provisioning

When a user logs in for the first time, create a local LMS user record from IDP
claims.

### Scheduled Registry Sync

Optional background sync from user registries, especially for users who have not
logged in yet.

### SCIM Support

Add SCIM 2.0 later if enterprise customers need automated user lifecycle
management.

## Data Model Outline

Core tables:

```text
tenants
idp_connections
users
roles
permissions
user_roles

courses
course_versions
course_modules
lessons
content_items
attachments

cohorts
enrollments
learning_paths
learning_path_courses

quizzes
questions
quiz_attempts
assignments
submissions
grades

progress_events
lesson_progress
course_progress
completion_records
certificates

notifications
notification_preferences
audit_logs
```

Important indexes:

```text
users(tenant_id, idp_issuer, external_subject)
courses(tenant_id, status)
enrollments(tenant_id, user_id, course_id)
progress_events(tenant_id, user_id, course_id, created_at)
audit_logs(tenant_id, actor_user_id, created_at)
```

## Frontend Structure

Suggested React structure:

```text
frontend/
  src/
    app/
    auth/
    api/
    routes/
    layouts/
    features/
      courses/
      enrollments/
      assessments/
      progress/
      admin/
      reports/
    components/
    permissions/
    hooks/
```

Main screens:

- Login redirect and auth callback
- Dashboard
- Course catalog
- My learning
- Course player
- Quiz and assessment screens
- Instructor course builder
- Enrollment management
- User and group management
- Reports
- Admin settings
- IDP configuration

Design the course player responsively from the start, because it will influence
the future mobile experience.

## Mobile Readiness

Before building mobile apps, make these architecture choices:

- Use OIDC Authorization Code Flow with PKCE
- Keep APIs mobile-friendly and decoupled from web page structure
- Avoid browser-only assumptions in auth and session handling
- Use short-lived access tokens
- Plan refresh-token handling for native apps
- Design notifications with future push tokens in mind
- Serve files and videos through signed URLs
- Version APIs
- Avoid exposing internal database IDs if public-safe identifiers may be needed

Future mobile options:

- React Native for code and team skill reuse with React
- Flutter for strong UI consistency
- Native iOS and Android if product requirements justify it

## Security Requirements

Minimum baseline:

- JWT signature validation using JWKS
- Issuer, audience, expiry, and claim validation
- Authorization enforced in the backend
- Tenant isolation on every query
- Audit logs for admin and grading actions
- CSRF protection if cookies are used
- Secure CORS configuration
- Rate limits on sensitive endpoints
- Object storage access through signed URLs
- Encryption at rest through managed database and storage services
- Secrets stored in a secret manager
- MFA delegated to the IDP
- No passwords stored in the LMS

For high-trust or regulated environments:

- SOC 2-ready audit trail
- Data retention policies
- Admin impersonation controls
- Consent and privacy controls
- FERPA or GDPR readiness depending on market
- Backup and restore testing

## Operational Plan

Infrastructure:

- Containerized Go API
- Containerized React frontend or static build served through a CDN
- PostgreSQL
- Redis
- Object storage
- Background worker
- Queue
- Observability stack

CI/CD:

- Backend tests
- Frontend tests
- Database migration checks
- OpenAPI validation
- Linting
- Security scanning
- Docker image builds
- Environment-based deployments

Observability:

- Structured logs
- Request tracing
- Metrics
- Error tracking
- Audit logs
- Business events

## Delivery Phases

### Phase 1: Foundation

- Go backend skeleton
- React frontend skeleton
- PostgreSQL migrations
- Tenant model
- OIDC login
- Local user provisioning
- Basic RBAC
- `/me` endpoint
- Admin user bootstrap

### Phase 2: Core LMS

- Course CRUD
- Module and lesson builder
- File upload
- Course catalog
- Enrollment
- Learner dashboard
- Course player
- Progress tracking

### Phase 3: Assessments

- Quiz builder
- Question management
- Attempts
- Grading
- Completion rules
- Certificates

### Phase 4: Admin And Reporting

- User management
- Cohorts and groups
- Bulk enrollment
- Progress reports
- Completion reports
- Audit logs
- Notification system

### Phase 5: Enterprise Integrations

- Multiple IDP connections
- IDP group mapping
- Scheduled registry sync
- SCIM provisioning
- Webhooks
- External reporting exports

### Phase 6: Mobile Apps

- API hardening for mobile
- Push notifications
- Offline-friendly course metadata
- Mobile course player
- Native auth flow
- Download controls if required

## Important Early Decisions

Resolve these before implementation:

1. Is this single-tenant or multi-tenant from day one?
2. Which IDP should be supported first?
3. Will users self-enroll, be assigned courses, or both?
4. Do courses need versioning?
5. Do you need SCORM or xAPI compatibility?
6. Are assessments simple quizzes only, or do you need assignments and grading workflows?
7. Do you need certificates?
8. What user registries must be integrated first?
9. Will content include video hosting or only embedded and external video?
10. What compliance requirements apply?

## Pragmatic MVP Scope

A solid first version should include:

- OIDC login
- Tenant-aware backend
- Local LMS user projection
- Admin, instructor, and learner roles
- Course creation
- Lessons with text, file, and video link support
- Enrollment
- Learner dashboard
- Course player
- Progress tracking
- Simple quizzes
- Completion status
- Basic reports

This gives the project a real LMS foundation without overbuilding the enterprise
integration layer too early.
