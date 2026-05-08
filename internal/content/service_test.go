package content

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"toucan/internal/sections"
)

type mockSectionLookup struct {
	mock.Mock
}

func (m *mockSectionLookup) Get(ctx context.Context, id string) (sections.Section, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sections.Section), args.Error(1)
}

func TestServiceCreateContentItem(t *testing.T) {
	mockRepo := new(MockRepository)
	mockSections := new(mockSectionLookup)
	service := NewService(mockRepo, mockSections)

	sectionID := "sec-123"
	ctx := context.Background()

	mockSections.On("Get", ctx, sectionID).Return(sections.Section{ID: sectionID}, nil)
	mockRepo.On("Create", mock.AnythingOfType("content.Item")).Return(func(item Item) Item {
		item.ID = "content-1"
		return item
	}, nil)

	item, err := service.Create(ctx, CreateItemInput{
		SectionID: sectionID,
		Title:     "Welcome Video",
		Summary:   "Overview of the platform.",
		Type:      TypeVideo,
		Position:  1,
	})

	assert.NoError(t, err)
	assert.Equal(t, "content-1", item.ID)
	assert.Equal(t, sectionID, item.SectionID)
	mockRepo.AssertExpectations(t)
	mockSections.AssertExpectations(t)
}
