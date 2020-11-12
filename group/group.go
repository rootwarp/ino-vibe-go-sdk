package group

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
	"github.com/rootwarp/ino-vibe-go-sdk/user"
)

const (
	permitGroupTreeDepth = 10
)

var (
	serverURL = "grpc.ino-vibe.ino-on.dev:443"

	// ErrGroupNonExist describes requested group is not exist on system.
	ErrGroupNonExist = errors.New("Group does not exist")
)

// Group is data structure for describing Auth0 Group.
type Group struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Children []Group `json:"children"`
}

type groupNode struct {
	ID       string
	Name     string
	Parent   *groupNode
	Children []*groupNode
}

// Client is client for Group.
type Client interface {
	List(ctx context.Context, groupID string) ([]Group, error)
	GetName(ctx context.Context, groupID string) (string, error)
	GetID(ctx context.Context, groupName string) (string, error)
	GetIDs(ctx context.Context, groupName []string) ([]string, error)
	GetChildGroups(ctx context.Context, groupID string) ([]Group, error)
	GetParentUsers(ctx context.Context, groupID string) ([]string, error)
	GetMembers(ctx context.Context, groupID string) ([]user.User, error)
}

type client struct {
	oauthToken  *oauth2.Token
	groupClient pb.GroupServiceClient
}

func (c *client) getGroupClient() pb.GroupServiceClient {
	if c.groupClient == nil {
		if c.oauthToken == nil {
			log.Panicln(errors.New("No credentials"))
		}

		certPool, err := x509.SystemCertPool()
		if err != nil {
			log.Panicln(err)
		}

		creds := credentials.NewClientTLSFromCert(certPool, "")
		conn, _ := grpc.Dial(
			serverURL,
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(oauth.NewOauthAccess(c.oauthToken)),
		)
		c.groupClient = pb.NewGroupServiceClient(conn)
	}

	return c.groupClient
}

func (c *client) List(ctx context.Context, groupID string) ([]Group, error) {
	cli := c.getGroupClient()

	listCli, err := cli.List(ctx, &pb.GroupRequest{Groupid: groupID})
	if err != nil {
		return []Group{}, nil
	}

	pbGroups := map[string]*pb.Group{}
	groupNodes := map[string]*groupNode{}
	rootNodes := make([]*groupNode, 0)

	for {
		group, err := listCli.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return []Group{}, err
			}
		}

		pbGroups[group.Groupid] = group
		groupNodes[group.Groupid] = &groupNode{ID: group.Groupid, Name: group.Name, Children: []*groupNode{}}

		if group.ParentId == "" {
			rootNodes = append(rootNodes, groupNodes[group.Groupid])
		}
	}

	for _, g := range pbGroups {
		if parent, ok := groupNodes[g.ParentId]; ok {
			parent.Children = append(parent.Children, groupNodes[g.Groupid])
		}
	}

	return c.traverse(rootNodes, 0), nil
}

func (c *client) printTree(roots []*groupNode, depth int) {
	if depth > permitGroupTreeDepth {
		return
	}

	for _, g := range roots {
		for i := 0; i < depth; i++ {
			fmt.Printf("-")
		}

		fmt.Println(g.Name)
		c.printTree(g.Children, depth+1)
	}
}

func (c *client) traverse(roots []*groupNode, depth int) []Group {
	if depth > permitGroupTreeDepth {
		return []Group{}
	}

	retGroups := make([]Group, len(roots))
	for i, g := range roots {
		retGroups[i] = Group{ID: g.ID, Name: g.Name}
		retGroups[i].Children = c.traverse(g.Children, depth+1)

	}

	return retGroups
}

// GetName returns group's name.
func (c *client) GetName(ctx context.Context, groupID string) (string, error) {
	cli := c.getGroupClient()

	resp, err := cli.Detail(ctx, &pb.GroupRequest{Groupid: groupID})
	if err != nil {
		log.Println(err)
		return "", err
	}

	if resp.ResultCode == pb.ResponseCode_NON_EXIST {
		return "", ErrGroupNonExist
	}

	return resp.Groups[0].GetName(), nil
}

// GetID returns group's name.
func (c *client) GetID(ctx context.Context, groupName string) (string, error) {
	cli := c.getGroupClient()

	resp, err := cli.FindByID(ctx, &pb.GroupFindRequest{Names: []string{groupName}})
	if err != nil {
		log.Println(err)
		return "", err
	}

	if len(resp.Groups) == 0 {
		return "", ErrGroupNonExist
	}

	return resp.Groups[0].GetGroupid(), nil
}

// GetIDs returns slice of group name.
func (c *client) GetIDs(ctx context.Context, groupNames []string) ([]string, error) {
	cli := c.getGroupClient()

	resp, err := cli.FindByID(ctx, &pb.GroupFindRequest{Names: groupNames})
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	if len(resp.Groups) == 0 {
		return []string{}, ErrGroupNonExist
	}

	groupIDs := make([]string, len(resp.Groups))
	for i, group := range resp.Groups {
		groupIDs[i] = group.GetGroupid()
	}

	return groupIDs, nil
}

// GetChildGroups returns tree based child groups.
func (c *client) GetChildGroups(ctx context.Context, groupID string) ([]Group, error) {
	cli := c.getGroupClient()

	resp, err := cli.Childs(ctx, &pb.GroupRequest{Groupid: groupID})
	if err != nil {
		return []Group{}, err
	}

	if resp.ResultCode == pb.ResponseCode_NON_EXIST {
		return []Group{}, ErrGroupNonExist
	}

	groups := make([]Group, len(resp.Groups))
	for i, group := range resp.Groups {
		groups[i] = Group{Name: group.Name, ID: group.Groupid}
	}

	return groups, nil
}

// GetParentUsers return list of all users in parent groups.
// Return value of []string contains email addresses of users.
func (c *client) GetParentUsers(ctx context.Context, groupID string) ([]string, error) {
	cli := c.getGroupClient()

	resp, err := cli.ParentUsers(ctx, &pb.GroupRequest{Groupid: groupID})
	if err != nil {
		return []string{}, err
	}

	if resp.ResponseCode == pb.ResponseCode_NON_EXIST {
		return []string{}, ErrGroupNonExist
	}

	emails := make([]string, len(resp.Emails))
	for i, email := range resp.Emails {
		emails[i] = email
	}

	return emails, nil
}

// GetMembers returns list of users who joined selected group.
func (c *client) GetMembers(ctx context.Context, groupID string) ([]user.User, error) {
	cli := c.getGroupClient()

	memberResp, err := cli.Members(ctx, &pb.GroupRequest{Groupid: groupID})
	if err != nil {
		return []user.User{}, err
	}

	respUsers := make([]user.User, len(memberResp.Users))
	for i, pbUser := range memberResp.Users {
		respUsers[i] = user.User{UserID: pbUser.UserId, Email: pbUser.Email}
	}

	return respUsers, nil
}

// NewClient creates new group client.
func NewClient() (Client, error) {
	token, err := iv_auth.LoadCredentials()
	if err != nil {
		log.Panicln(err)
	}

	return &client{oauthToken: token}, nil
}
