package retl_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// mockService implements retl.Service for unit tests by recording calls via
// testify's mock.Mock.
type mockService struct {
	mock.Mock
}

func (m *mockService) CreateRetlSource(ctx context.Context, source *retl.RETLSourceCreateRequest) (*retl.RETLSource, error) {
	args := m.Called(ctx, source)
	if v := args.Get(0); v != nil {
		return v.(*retl.RETLSource), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockService) GetRetlSource(ctx context.Context, id string) (*retl.RETLSource, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*retl.RETLSource), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockService) UpdateRetlSource(ctx context.Context, id string, source *retl.RETLSourceUpdateRequest) (*retl.RETLSource, error) {
	args := m.Called(ctx, id, source)
	if v := args.Get(0); v != nil {
		return v.(*retl.RETLSource), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockService) DeleteRetlSource(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockService) CreateConnection(ctx context.Context, req *retl.CreateRETLConnectionRequest) (*retl.RETLConnection, error) {
	args := m.Called(ctx, req)
	if v := args.Get(0); v != nil {
		return v.(*retl.RETLConnection), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockService) GetConnection(ctx context.Context, id string) (*retl.RETLConnection, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*retl.RETLConnection), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockService) UpdateConnection(ctx context.Context, id string, req *retl.UpdateRETLConnectionRequest) (*retl.RETLConnection, error) {
	args := m.Called(ctx, id, req)
	if v := args.Get(0); v != nil {
		return v.(*retl.RETLConnection), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockService) DeleteConnection(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
