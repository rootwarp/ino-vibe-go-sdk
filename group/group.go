package group

import (
	"context"
	"crypto/x509"
	"errors"
	"log"

	pb "bitbucket.org/ino-on/ino-vibe-api"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	iv_auth "github.com/rootwarp/ino-vibe-go-sdk/auth"
)

// Reader provides retrieve operations for Group.
type Reader interface {
	GetName(ctx context.Context, groupID string) (string, error)
	GetID(ctx context.Context, groupName string) (string, error)
	GetIDs(ctx context.Context, groupName []string) (string, error)
	GetChildGroups(ctx context.Context, groupID string) ([]Group, error)
	GetParentUsers(ctx context.Context, groupID string) ([]string, error)
}

// Group is data structure for describing Auth0 Group.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Writer provides update operations for Group.
type Writer interface{}

const (
	serverURL = "group.ino-vibe.ino-on.dev:443"
)

var (
	// ErrGroupNonExist describes requested group is not exist on system.
	ErrGroupNonExist = errors.New("Group does not exist")
)

// Client is client for Group.
type Client struct {
	oauthToken  *oauth2.Token
	groupClient pb.GroupServiceClient
}

func (c *Client) getGroupClient() pb.GroupServiceClient {
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

// GetName returns group's name.
func (c *Client) GetName(ctx context.Context, groupID string) (string, error) {
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
func (c *Client) GetID(ctx context.Context, groupName string) (string, error) {
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
func (c *Client) GetIDs(ctx context.Context, groupNames []string) ([]string, error) {
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
func (c *Client) GetChildGroups(ctx context.Context, groupID string) ([]Group, error) {
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
func (c *Client) GetParentUsers(ctx context.Context, groupID string) ([]string, error) {
	cli := c.getGroupClient()

	resp, err := cli.NestedUsers(ctx, &pb.GroupRequest{Groupid: groupID})
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

// NewClient creates new group client.
func NewClient() (*Client, error) {
	token, err := iv_auth.LoadCredentials()
	if err != nil {
		log.Panicln(err)
	}

	return &Client{oauthToken: token}, nil
}
