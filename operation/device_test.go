package operation

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	"github.com/rootwarp/ino-vibe-go-sdk/device"
	"github.com/rootwarp/ino-vibe-go-sdk/group"
)

var (
	groupClient  group.Client
	deviceClient device.Client
)

func init() {
	var err error

	groupClient, err = group.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	deviceClient, err = device.NewClient()
	if err != nil {
		log.Fatal(err)
	}
}

type List struct {
	items  []interface{}
	marker map[interface{}]bool
}

func (l *List) Append(item interface{}) {
	if l.items == nil {
		l.items = []interface{}{}
		l.marker = map[interface{}]bool{}
	}

	l.items = append(l.items, item)
	l.marker[item] = true
}

func (l *List) Contain(item interface{}) bool {
	_, ok := l.marker[item]
	return ok
}

func (l *List) Print() {
	if l.items == nil {
		return
	}

	for _, item := range l.items {
		fmt.Printf("%+v\n", item)
	}
}

func TestWorkingList(t *testing.T) {
	const groupName = "코원 에너지"

	// Get Group ID.
	ctx := context.Background()
	groupID, _ := groupClient.GetID(ctx, groupName)

	groupList := List{}

	children, _ := groupClient.GetChildGroups(ctx, groupID)
	for _, group := range children {
		groupList.Append(group.ID)
	}

	deviceResp, _ := deviceClient.List(ctx, pb.InstallStatus_Installed)
	filteredDevice := make([]*pb.Device, 0)

	// Filter
	for _, device := range deviceResp.Devices {
		if groupList.Contain(device.GroupId) {
			filteredDevice = append(filteredDevice, device)
		}
	}

	for _, device := range filteredDevice {
		groupName, _ := groupClient.GetName(ctx, device.GroupId)
		installDate := time.Unix(device.InstallDate.Seconds, 0)
		fmt.Printf("%s, %s, %s, %s, %s, %d\n", device.Devid, device.Alias, groupName, device.Installer, installDate, device.Battery)
	}
}
