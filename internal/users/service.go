package users

import (
	"context"
	"fmt"
	"strings"
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

func (s *service) Get(ctx context.Context, id string) (User, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetByExternalSubject(ctx context.Context, subject string) (User, error) {
	return s.repo.GetByExternalSubject(ctx, subject)
}

func (s *service) List(ctx context.Context, filter ListFilter) (ListResult, error) {
	return s.repo.List(ctx, filter)
}
