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
		serverURL = target + "-" + serverURL
	}
	fmt.Println(serverURL)
}

func TestGroupList(t *testing.T) {
	partialRootGroups := []string{
		"0bee7b43-0b57-4b54-9062-430e2bd3fa79", // Ino-on
		"406b3434-7ddf-4af8-b357-cc144415bcb7", // SK E&S
		"1590fe0a-e416-48f7-b9c7-d8f4f37f4d64", // testing
	}

	cli, _ := NewClient()
	ctx := context.Background()

	groups, _ := cli.List(ctx, "")
	groupIDs := map[string]bool{}
	for _, g := range groups {
		groupIDs[g.ID] = true
		assert.Nil(t, g.Parent)
	}

	for _, rootID := range partialRootGroups {
		assert.Contains(t, groupIDs, rootID)
	}
}

func TestGroupListForSelected(t *testing.T) {
	/*
		이노온 contains
		이노온 - 개발
		이노온 - 운영
	*/

	cli, _ := NewClient()
	ctx := context.Background()

	groups, _ := cli.List(ctx, "0bee7b43-0b57-4b54-9062-430e2bd3fa79")

	assert.Equal(t, 1, len(groups))
	assert.Equal(t, "0bee7b43-0b57-4b54-9062-430e2bd3fa79", groups[0].ID)

	childrenIDs := make([]string, len(groups[0].Children))
	for i, children := range groups[0].Children {
		childrenIDs[i] = children.ID
	}

	assert.Contains(t, childrenIDs, "607f9db4-7eee-4a08-894d-356c8a462ae1")
	assert.Contains(t, childrenIDs, "b09a8694-6ccb-4cb7-9ffa-57681869f54d")
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
	emails, err := cli.GetParentUsers(ctx, groupID)

	assert.Contains(t, emails, "child_tester@ino-on.com")
	assert.Contains(t, emails, "parent_tester@ino-on.com")

	groupID, err = cli.GetID(ctx, "이노온")
	emails, err = cli.GetParentUsers(ctx, groupID)

	assert.Nil(t, err)
	assert.Contains(t, emails, "parent_tester@ino-on.com")
	assert.NotContains(t, emails, "child_tester@ino-on.com")
}

func TestGroupMembers(t *testing.T) {
	cli, _ := NewClient()
	ctx := context.Background()

	users, err := cli.GetMembers(ctx, "0bee7b43-0b57-4b54-9062-430e2bd3fa79")

	assert.Nil(t, err)

	userEmailMap := map[string]bool{}
	for _, user := range users {
		userEmailMap[user.Email] = true
	}

	assert.Contains(t, userEmailMap, "ino-vibe@ino-on.com")
	assert.Contains(t, userEmailMap, "develop@ino-on.com")
}

func TestGroupMembersWithEmptyGroupID(t *testing.T) {
	cli, _ := NewClient()
	ctx := context.Background()

	users, err := cli.GetMembers(ctx, "")

	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))
}
