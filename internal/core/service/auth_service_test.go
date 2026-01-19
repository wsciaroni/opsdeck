package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// MockUserRepository is a mock implementation of port.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockOrganizationRepository is a mock implementation of port.OrganizationRepository
type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) AddMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role string) error {
	args := m.Called(ctx, orgID, userID, role)
	return args.Error(0)
}

func (m *MockOrganizationRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserMembership, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserMembership), args.Error(1)
}

// MockOIDCProvider is a mock implementation of port.OIDCProvider
type MockOIDCProvider struct {
	mock.Mock
}

func (m *MockOIDCProvider) AuthCodeURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockOIDCProvider) ExchangeCode(ctx context.Context, code string) (*port.UserInfo, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.UserInfo), args.Error(1)
}

func TestGetLoginURL(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockOIDC := new(MockOIDCProvider)
	logger := slog.Default()
	service := NewAuthService(mockRepo, mockOrgRepo, mockOIDC, logger)

	expectedURL := "https://accounts.google.com/o/oauth2/auth?state=state-random-string"
	mockOIDC.On("AuthCodeURL", "state-random-string").Return(expectedURL)

	url := service.GetLoginURL()
	assert.Equal(t, expectedURL, url)
	mockOIDC.AssertExpectations(t)
}

func TestLoginFromProvider_NewUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockOIDC := new(MockOIDCProvider)
	logger := slog.Default()
	service := NewAuthService(mockRepo, mockOrgRepo, mockOIDC, logger)

	ctx := context.Background()
	code := "test-code"
	userInfo := &port.UserInfo{
		Email:     "test@example.com",
		Name:      "Test User",
		AvatarURL: "http://example.com/avatar.jpg",
	}

	mockOIDC.On("ExchangeCode", ctx, code).Return(userInfo, nil)
	mockRepo.On("GetByEmail", ctx, userInfo.Email).Return(nil, nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == userInfo.Email && u.Name == userInfo.Name && u.AvatarURL == userInfo.AvatarURL && u.Role == domain.RolePublic
	})).Return(nil)

	// Expect Org Creation
	mockOrgRepo.On("Create", ctx, mock.MatchedBy(func(o *domain.Organization) bool {
		return o.Name == "Personal Workspace"
	})).Return(nil)

	// Expect Add Member
	mockOrgRepo.On("AddMember", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"), "owner").Return(nil)

	user, err := service.LoginFromProvider(ctx, code)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userInfo.Email, user.Email)
	assert.Equal(t, userInfo.Name, user.Name)

	mockOIDC.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockOrgRepo.AssertExpectations(t)
}

func TestLoginFromProvider_ExistingUser_NoUpdate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockOIDC := new(MockOIDCProvider)
	logger := slog.Default()
	service := NewAuthService(mockRepo, mockOrgRepo, mockOIDC, logger)

	ctx := context.Background()
	code := "test-code"
	userInfo := &port.UserInfo{
		Email:     "test@example.com",
		Name:      "Test User",
		AvatarURL: "http://example.com/avatar.jpg",
	}
	existingUser := &domain.User{
		ID:        uuid.New(),
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.AvatarURL,
		Role:      domain.RolePublic,
	}

	mockOIDC.On("ExchangeCode", ctx, code).Return(userInfo, nil)
	mockRepo.On("GetByEmail", ctx, userInfo.Email).Return(existingUser, nil)

	user, err := service.LoginFromProvider(ctx, code)
	assert.NoError(t, err)
	assert.Equal(t, existingUser, user)

	mockOIDC.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestLoginFromProvider_ExistingUser_WithUpdate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockOIDC := new(MockOIDCProvider)
	logger := slog.Default()
	service := NewAuthService(mockRepo, mockOrgRepo, mockOIDC, logger)

	ctx := context.Background()
	code := "test-code"
	userInfo := &port.UserInfo{
		Email:     "test@example.com",
		Name:      "New Name",
		AvatarURL: "http://example.com/new-avatar.jpg",
	}
	existingUser := &domain.User{
		ID:        uuid.New(),
		Email:     userInfo.Email,
		Name:      "Old Name",
		AvatarURL: "http://example.com/old-avatar.jpg",
		Role:      domain.RolePublic,
	}

	mockOIDC.On("ExchangeCode", ctx, code).Return(userInfo, nil)
	mockRepo.On("GetByEmail", ctx, userInfo.Email).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.Name == userInfo.Name && u.AvatarURL == userInfo.AvatarURL
	})).Return(nil)

	user, err := service.LoginFromProvider(ctx, code)
	assert.NoError(t, err)
	assert.Equal(t, userInfo.Name, user.Name)
	assert.Equal(t, userInfo.AvatarURL, user.AvatarURL)

	mockOIDC.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateSession(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockOIDC := new(MockOIDCProvider)
	logger := slog.Default()
	service := NewAuthService(mockRepo, mockOrgRepo, mockOIDC, logger)

	ctx := context.Background()
	user := &domain.User{
		ID: uuid.New(),
	}

	token, err := service.CreateSession(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user.ID.String(), token)
}
