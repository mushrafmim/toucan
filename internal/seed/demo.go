package seed

import (
	"context"
	"toucan/internal/content"
	"toucan/internal/courses"
	"toucan/internal/identity"
	"toucan/internal/sections"
	"toucan/internal/users"
)

func Demo(
	userService users.Service,
	courseService *courses.Service,
	sectionService *sections.Service,
	contentService *content.Service,
) {
	ctx := context.Background()

	// Ensure a seed user exists to act as the creator
	const seedSubject = "system-seed"
	user, err := userService.GetByExternalSubject(ctx, seedSubject)
	if err != nil {
		user, err = userService.Create(ctx, users.CreateUserRequest{
			ExternalSubject: seedSubject,
			Email:           "seed@toucan.local",
			DisplayName:     "Seed Robot",
			Roles:           []users.Role{users.RoleAdmin},
		})
		if err != nil {
			return
		}
	}

	// Create an authorized context
	ctx = identity.ContextWithPrincipal(ctx, identity.Principal{
		Subject: user.ExternalSubject,
		Email:   user.Email,
		Roles:   []string{"admin"},
	})

	goCourse, err := courseService.Create(ctx, courses.CreateCourseInput{
		Title:       "Go for LMS Services",
		Summary:     "Build and maintain backend services for the learning platform.",
		Description: "Covers HTTP handlers, modular domain design, in-memory testing, and API composition patterns.",
		Category:    "Engineering",
		Level:       courses.LevelIntermediate,
		Tags:        []string{"go", "backend", "api"},
	})
	if err != nil {
		return
	}

	fundamentals, err := sectionService.Create(ctx, sections.CreateSectionInput{
		CourseID: goCourse.ID,
		Title:    "Foundations",
		Summary:  "Start with the service architecture and routing model.",
		Position: 1,
	})
	if err != nil {
		return
	}

	implementation, err := sectionService.Create(ctx, sections.CreateSectionInput{
		CourseID: goCourse.ID,
		Title:    "Implementation",
		Summary:  "Wire handlers, repositories, and frontend consumers together.",
		Position: 2,
	})
	if err != nil {
		return
	}

	_, _ = contentService.Create(ctx, content.CreateItemInput{
		SectionID: fundamentals.ID,
		Title:     "Architecture Walkthrough",
		Summary:   "Overview of the modular LMS backend shape.",
		Type:      content.TypeVideo,
		Position:  1,
		SourceURL: "https://cdn.example.test/architecture-walkthrough.mp4",
		Metadata:  map[string]any{"duration_minutes": 14},
	})

	_, _ = contentService.Create(ctx, content.CreateItemInput{
		SectionID: fundamentals.ID,
		Title:     "Service Boundaries",
		Summary:   "Reference notes for courses, sections, and content domains.",
		Type:      content.TypePDF,
		Position:  2,
		SourceURL: "https://cdn.example.test/service-boundaries.pdf",
	})

	_, _ = contentService.Create(ctx, content.CreateItemInput{
		SectionID: implementation.ID,
		Title:     "Course Detail Aggregation",
		Summary:   "How the frontend composes course pages from multiple APIs.",
		Type:      content.TypeRichText,
		Position:  1,
		Body:      "Fetch the course first, then its sections, then content for each section.",
	})

	opsCourse, err := courseService.Create(ctx, courses.CreateCourseInput{
		Title:       "Operator Onboarding",
		Summary:     "Operational basics for running and supporting Toucan.",
		Description: "Introduces deployment workflows, debugging patterns, and support expectations for the LMS.",
		Category:    "Operations",
		Level:       courses.LevelBeginner,
		Tags:        []string{"ops", "support"},
	})
	if err != nil {
		return
	}

	opsSection, err := sectionService.Create(ctx, sections.CreateSectionInput{
		CourseID: opsCourse.ID,
		Title:    "First Day Setup",
		Summary:  "Environment and support tooling basics.",
		Position: 1,
	})
	if err != nil {
		return
	}

	_, _ = contentService.Create(ctx, content.CreateItemInput{
		SectionID: opsSection.ID,
		Title:     "Access Checklist",
		Summary:   "A short checklist for first-day setup tasks.",
		Type:      content.TypeLink,
		Position:  1,
		SourceURL: "https://intranet.example.test/toucan-access-checklist",
	})
}
