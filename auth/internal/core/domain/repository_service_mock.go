// Code generated by mockery v2.33.1. DO NOT EDIT.

package domain

import mock "github.com/stretchr/testify/mock"

// MockRepositoryService is an autogenerated mock type for the RepositoryService type
type MockRepositoryService struct {
	mock.Mock
}

// ApiKeyRepository provides a mock function with given fields:
func (_m *MockRepositoryService) ApiKeyRepository() ApiKeyRepository {
	ret := _m.Called()

	var r0 ApiKeyRepository
	if rf, ok := ret.Get(0).(func() ApiKeyRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ApiKeyRepository)
		}
	}

	return r0
}

// NewMockRepositoryService creates a new instance of MockRepositoryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepositoryService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepositoryService {
	mock := &MockRepositoryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}