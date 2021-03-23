// Code generated by mockery v1.0.0. DO NOT EDIT.

package group

import (
	context "context"

	user "github.com/rootwarp/ino-vibe-go-sdk/user"
	mock "github.com/stretchr/testify/mock"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, name, parent
func (_m *MockClient) Create(ctx context.Context, name string, parent *Group) (*Group, error) {
	ret := _m.Called(ctx, name, parent)

	var r0 *Group
	if rf, ok := ret.Get(0).(func(context.Context, string, *Group) *Group); ok {
		r0 = rf(ctx, name, parent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *Group) error); ok {
		r1 = rf(ctx, name, parent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, groupID
func (_m *MockClient) Delete(ctx context.Context, groupID string) error {
	ret := _m.Called(ctx, groupID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, groupID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetChildGroups provides a mock function with given fields: ctx, groupID
func (_m *MockClient) GetChildGroups(ctx context.Context, groupID string) ([]Group, error) {
	ret := _m.Called(ctx, groupID)

	var r0 []Group
	if rf, ok := ret.Get(0).(func(context.Context, string) []Group); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetID provides a mock function with given fields: ctx, groupName
func (_m *MockClient) GetID(ctx context.Context, groupName string) (string, error) {
	ret := _m.Called(ctx, groupName)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, groupName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetIDs provides a mock function with given fields: ctx, groupName
func (_m *MockClient) GetIDs(ctx context.Context, groupName []string) ([]string, error) {
	ret := _m.Called(ctx, groupName)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, []string) []string); ok {
		r0 = rf(ctx, groupName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, groupName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMembers provides a mock function with given fields: ctx, groupID
func (_m *MockClient) GetMembers(ctx context.Context, groupID string) ([]user.User, error) {
	ret := _m.Called(ctx, groupID)

	var r0 []user.User
	if rf, ok := ret.Get(0).(func(context.Context, string) []user.User); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetName provides a mock function with given fields: ctx, groupID
func (_m *MockClient) GetName(ctx context.Context, groupID string) (string, error) {
	ret := _m.Called(ctx, groupID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, groupID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetParentUsers provides a mock function with given fields: ctx, groupID
func (_m *MockClient) GetParentUsers(ctx context.Context, groupID string) ([]string, error) {
	ret := _m.Called(ctx, groupID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, groupID
func (_m *MockClient) List(ctx context.Context, groupID string) ([]Group, error) {
	ret := _m.Called(ctx, groupID)

	var r0 []Group
	if rf, ok := ret.Get(0).(func(context.Context, string) []Group); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
