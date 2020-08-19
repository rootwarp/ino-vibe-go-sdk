package group

import (
	"context"
	"fmt"
	"os"
	"testing"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureGroups = []pb.Group{
		{Groupid: "607f9db4-7eee-4a08-894d-356c8a462ae1", Name: "이노온-개발"},
		{Groupid: "b09a8694-6ccb-4cb7-9ffa-57681869f54d", Name: "이노온-운영"},
	}
)

func init() {
	target := os.Getenv("TEST_TARGET")
	if target != "" {
		serverURL = "grpc-dev.ino-vibe.ino-on.dev:443"

		if target == "feature" {
			serverURL = target + "-" + serverURL
		}
	}
	fmt.Println(serverURL)
}

func TestGroupGetName(t *testing.T) {
	expectations := []struct {
		GroupID    string
		ExpectName string
		ExpectErr  error
	}{
		{
			GroupID:    fixtureGroups[0].Groupid,
			ExpectName: fixtureGroups[0].Name,
			ExpectErr:  nil,
		},
		{
			GroupID:    "non-exist",
			ExpectName: "",
			ExpectErr:  ErrGroupNonExist,
		},
	}

	cli, _ := NewClient()
	ctx := context.Background()

	for _, expect := range expectations {
		name, err := cli.GetName(ctx, expect.GroupID)

		assert.Equal(t, expect.ExpectErr, err)
		assert.Equal(t, expect.ExpectName, name)
	}
}

func TestGroupGetIDOK(t *testing.T) {
	expectations := []struct {
		GroupName string
		ExpectID  string
		ExpectErr error
	}{
		{
			GroupName: fixtureGroups[0].Name,
			ExpectID:  fixtureGroups[0].Groupid,
			ExpectErr: nil,
		},
		{
			GroupName: "non-exist",
			ExpectID:  "",
			ExpectErr: ErrGroupNonExist,
		},
	}

	cli, _ := NewClient()
	ctx := context.Background()

	for _, expect := range expectations {
		groupID, err := cli.GetID(ctx, expect.GroupName)

		assert.Equal(t, expect.ExpectErr, err)
		assert.Equal(t, expect.ExpectID, groupID)
	}
}

func TestGroupGetIDsOK(t *testing.T) {
	cli, _ := NewClient()
	ctx := context.Background()

	groupIDs, err := cli.GetIDs(ctx, []string{fixtureGroups[0].Name, fixtureGroups[1].Name})

	assert.Nil(t, err)
	assert.Equal(t, fixtureGroups[0].Groupid, groupIDs[0])
	assert.Equal(t, fixtureGroups[1].Groupid, groupIDs[1])
}

func TestGroupGetChildsOK(t *testing.T) {
	expectations := []struct {
		GroupID   string
		ExpectErr error
	}{
		{
			GroupID:   fixtureGroups[0].Groupid,
			ExpectErr: nil,
		},
		{
			GroupID:   "non-exist",
			ExpectErr: ErrGroupNonExist,
		},
	}

	cli, _ := NewClient()
	ctx := context.Background()

	for _, expect := range expectations {
		_, err := cli.GetChildGroups(ctx, expect.GroupID)

		assert.Equal(t, expect.ExpectErr, err)
	}
}

func TestGroupGetChildUsersOK(t *testing.T) {
	expectations := []struct {
		GroupID   string
		ExpectErr error
	}{
		{
			GroupID:   fixtureGroups[0].Groupid,
			ExpectErr: nil,
		},
		{
			GroupID:   "non-exist",
			ExpectErr: ErrGroupNonExist,
		},
	}

	cli, _ := NewClient()
	ctx := context.Background()

	for _, expect := range expectations {
		_, err := cli.GetParentUsers(ctx, expect.GroupID)

		assert.Equal(t, expect.ExpectErr, err)
	}
}

func TestGroupParentUsers(t *testing.T) {
	cli, _ := NewClient()
	ctx := context.Background()

	groupID, err := cli.GetID(ctx, "이노온-개발-서버")
	fmt.Println(groupID, err)

	emails, err := cli.GetParentUsers(ctx, groupID)

	assert.Contains(t, emails, "child_tester@ino-on.com")
	assert.Contains(t, emails, "parent_tester@ino-on.com")

	groupID, err = cli.GetID(ctx, "이노온")
	emails, err = cli.GetParentUsers(ctx, groupID)

	assert.Contains(t, emails, "parent_tester@ino-on.com")
	assert.NotContains(t, emails, "child_tester@ino-on.com")
}
