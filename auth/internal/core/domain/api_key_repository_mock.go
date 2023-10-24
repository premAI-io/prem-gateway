// Code generated by mockery v2.33.1. DO NOT EDIT.

package domain

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockApiKeyRepository is an autogenerated mock type for the ApiKeyRepository type
type MockApiKeyRepository struct {
	mock.Mock
}

// CreateApiKey provides a mock function with given fields: ctx, key
func (_m *MockApiKeyRepository) CreateApiKey(ctx context.Context, key ApiKey) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ApiKey) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteApiKey provides a mock function with given fields: ctx, id
func (_m *MockApiKeyRepository) DeleteApiKey(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllApiKeys provides a mock function with given fields: ctx
func (_m *MockApiKeyRepository) GetAllApiKeys(ctx context.Context) ([]ApiKey, error) {
	ret := _m.Called(ctx)

	var r0 []ApiKey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]ApiKey, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []ApiKey); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ApiKey)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApiKey provides a mock function with given fields: ctx, id
func (_m *MockApiKeyRepository) GetApiKey(ctx context.Context, id string) (*ApiKey, error) {
	ret := _m.Called(ctx, id)

	var r0 *ApiKey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*ApiKey, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *ApiKey); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ApiKey)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetServiceApiKey provides a mock function with given fields: ctx, service
func (_m *MockApiKeyRepository) GetServiceApiKey(ctx context.Context, service string) (*ApiKey, error) {
	ret := _m.Called(ctx, service)

	var r0 *ApiKey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*ApiKey, error)); ok {
		return rf(ctx, service)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *ApiKey); ok {
		r0 = rf(ctx, service)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ApiKey)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, service)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockApiKeyRepository creates a new instance of MockApiKeyRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockApiKeyRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockApiKeyRepository {
	mock := &MockApiKeyRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}