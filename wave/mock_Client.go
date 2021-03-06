// Code generated by mockery v1.0.0. DO NOT EDIT.

package wave

import (
	context "context"

	inovibe_api_v3 "bitbucket.org/ino-on/ino-vibe-api"
	mock "github.com/stretchr/testify/mock"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockClient) Close() {
	_m.Called()
}

// Detail provides a mock function with given fields: ctx, req
func (_m *MockClient) Detail(ctx context.Context, req *inovibe_api_v3.WaveDetailRequest) (*inovibe_api_v3.WaveDetailResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *inovibe_api_v3.WaveDetailResponse
	if rf, ok := ret.Get(0).(func(context.Context, *inovibe_api_v3.WaveDetailRequest) *inovibe_api_v3.WaveDetailResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*inovibe_api_v3.WaveDetailResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *inovibe_api_v3.WaveDetailRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, devid, offset, maxCount
func (_m *MockClient) List(ctx context.Context, devid string, offset int, maxCount int) ([]*inovibe_api_v3.WaveDetailItem, error) {
	ret := _m.Called(ctx, devid, offset, maxCount)

	var r0 []*inovibe_api_v3.WaveDetailItem
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*inovibe_api_v3.WaveDetailItem); ok {
		r0 = rf(ctx, devid, offset, maxCount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*inovibe_api_v3.WaveDetailItem)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, devid, offset, maxCount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
