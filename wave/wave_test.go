package wave

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	"github.com/stretchr/testify/assert"
)

func init() {
	target := os.Getenv("TEST_TARGET")
	if target != "" {
		serverURL = target + "-" + serverURL
	}
	fmt.Println(serverURL)
}

func TestWaveDetailSuccess(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		DeviceID    string
		WaveID      string
		GroupID     string
		Resolution  float32
		IntervalMs  uint32
		WaveSession int
		Notify      bool
		Created     time.Time
		X           []int32
		Y           []int32
		Z           []int32
	}{
		{
			DeviceID:    "00000125d02544fffefe108a",
			Created:     time.Date(2020, time.May, 26, 13, 51, 11, 0, time.FixedZone("JST", 9*3600)),
			WaveID:      "00000125d02544fffefe108a-1590468671",
			GroupID:     "0bee7b43-0b57-4b54-9062-430e2bd3fa79",
			Resolution:  0.244,
			IntervalMs:  10,
			WaveSession: 15,
			Notify:      false,
			Z: []int32{
				2524, 4888, 4976, 3224, 4828, 4936, 2496, 4544,
				4876, 2796, 3796, 4872, 3876, 2764, 4912, 4768,
				2496, 4872, 4984, 3288, 4836, 4920, 2468, 4484,
				4908, 2844, 3768, 4888, 3964, 2700, 4920, 4800,
			},
		},
		{
			DeviceID:    "00000125d02544fffefe15a6",
			Created:     time.Date(2020, time.May, 24, 9, 51, 10, 0, time.FixedZone("JST", 9*3600)),
			WaveID:      "00000125d02544fffefe15a6-1590281470",
			GroupID:     "0bee7b43-0b57-4b54-9062-430e2bd3fa79",
			Resolution:  0.244,
			IntervalMs:  10,
			WaveSession: 2,
			Notify:      false,
			Z: []int32{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				4636, 3452, 3412, 4404, 4400, 3408, 4392, 4436,
				3412, 4372, 4456, 3404, 4340, 4468, 3392, 3532,
				4584, 4100, 3508, 4560, 4136, 3472, 4568, 4180,
				3480, 4520, 4208, 3444, 3700, 4660, 3768, 3668,
			},
		},
	}

	c, _ := NewClient()
	for _, test := range tests {
		req := &pb.WaveDetailRequest{Waveid: test.WaveID}
		resp, err := c.Detail(ctx, req)

		assert.Nil(t, err)
		assert.Equal(t, test.DeviceID, resp.Wave.Devid)
		assert.Equal(t, test.WaveID, resp.Wave.Waveid)
		assert.Equal(t, test.GroupID, resp.Wave.Groupid)
		assert.Equal(t, test.Resolution, resp.Wave.Resolution)
		assert.Equal(t, test.IntervalMs, resp.Wave.Interval)
		assert.Equal(t, test.Notify, resp.Wave.Notify)
		assert.Equal(t, test.X, resp.Wave.X)
		assert.Equal(t, test.Y, resp.Wave.Y)
		assert.Equal(t, test.Z, resp.Wave.Z)
		assert.Equal(t, test.Created.Unix(), resp.Wave.Created.Seconds)
	}
}

func TestWaveDetailNonExist(t *testing.T) {
	c, _ := NewClient()

	ctx := context.Background()
	req := &pb.WaveDetailRequest{Waveid: "non-exist-wave-id"}
	resp, err := c.Detail(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, pb.ResponseCode_NON_EXIST, resp.ResponseCode)
}
