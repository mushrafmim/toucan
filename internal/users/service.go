package users

import (
	"context"
	"fmt"
	"strings"

	"toucan/internal/identity"
)

type service struct {
	repo *Repository
}

func NewService(repo *Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateUserRequest) (User, error) {
	if strings.TrimSpace(req.ExternalSubject) == "" || strings.TrimSpace(req.Email) == "" {
		return User{}, fmt.Errorf("external subject and email are required")
	}

	user := User{
		ExternalSubject: strings.TrimSpace(req.ExternalSubject),
		Email:           strings.TrimSpace(req.Email),
		DisplayName:     strings.TrimSpace(req.DisplayName),
		Roles:           req.Roles,
	}

	if len(user.Roles) == 0 {
		user.Roles = []Role{RoleLearner}
	}

	return s.repo.Create(ctx, user)
}

func (s *service) Update(ctx context.Context, u User) (User, error) {
	return s.repo.Update(ctx, u)
}

func (s *service) EnsureUser(ctx context.Context, principal identity.Principal) (User, error) {
	// 1. Map IDP roles to internal roles (strictly admin, instructor, learner)
	var internalRoles []Role
	for _, r := range principal.Roles {
		switch strings.ToLower(r) {
		case "admin":
			internalRoles = append(internalRoles, RoleAdmin)
		case "instructor":
			internalRoles = append(internalRoles, RoleInstructor)
		case "learner":
			internalRoles = append(internalRoles, RoleLearner)
		}
	}

	// 2. Default to learner if no valid roles found
	if len(internalRoles) == 0 {
		internalRoles = []Role{RoleLearner}
	}

	// 3. Check if user already exists
	existing, err := s.repo.GetByExternalSubject(ctx, principal.Subject)
	if err == nil {
		// Update roles and display name if they've changed in the IDP
		if !equalRoles(existing.Roles, internalRoles) || existing.DisplayName != principal.Name || existing.Email != principal.Email {
			existing.Roles = internalRoles
			existing.DisplayName = principal.Name
			existing.Email = principal.Email
			return s.repo.Update(ctx, existing)
		}
		return existing, nil
	}

	// 4. Not found, provision new user
	newUser := User{
		ExternalSubject: principal.Subject,
		Email:           principal.Email,
		DisplayName:     principal.Name,
		Roles:           internalRoles,
	}

	return s.repo.Create(ctx, newUser)
}

func (s *service) Get(ctx context.Context, id string) (User, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetByExternalSubject(ctx context.Context, subject string) (User, error) {
	return s.repo.GetByExternalSubject(ctx, subject)
}

func (s *service) List(ctx context.Context, filter ListFilter) (ListResult, error) {
	return s.repo.List(ctx, filter)
}

func equalRoles(a, b []Role) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[Role]int)
	for _, v := range a {
		m[v]++
	}
	for _, v := range b {
		if m[v] == 0 {
			return false
		}
		m[v]--
	}
	return true
}
