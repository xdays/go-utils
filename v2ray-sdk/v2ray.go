package v2ray

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	boundService "v2ray.com/core/app/proxyman/command"
	statsService "v2ray.com/core/app/stats/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

// User v2ray inbond user
type User struct {
	Level   uint32
	Email   string
	UUID    string
	AlterID uint32
}

// Client v2ray client wrapper
type Client struct {
	Host string
	Port int
}

// GetConnection initialize v2ray client grpc connection
func (v *Client) GetConnection() (*grpc.ClientConn, error) {
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", v.Host, v.Port), grpc.WithInsecure())
	return cmdConn, err
}

// GetBondClient get bond service client
func (v *Client) GetBondClient(g *grpc.ClientConn) boundService.HandlerServiceClient {
	return boundService.NewHandlerServiceClient(g)
}

// GetStatClient get state service client
func (v *Client) GetStatClient(g *grpc.ClientConn) statsService.StatsServiceClient {
	return statsService.NewStatsServiceClient(g)
}

// AddUser add user to inbond
func (v *Client) AddUser(c boundService.HandlerServiceClient, i string, u User) (*boundService.AlterInboundResponse, error) {
	resp, err := c.AlterInbound(context.Background(), &boundService.AlterInboundRequest{
		Tag: i,
		Operation: serial.ToTypedMessage(&boundService.AddUserOperation{
			User: &protocol.User{
				Level: u.Level,
				Email: u.Email,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               u.UUID,
					AlterId:          u.AlterID,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})
	return resp, err
}

// RemoveUser remove user from inbond
func (v *Client) RemoveUser(c boundService.HandlerServiceClient, i string, u User) (*boundService.AlterInboundResponse, error) {
	resp, err := c.AlterInbound(context.Background(), &boundService.AlterInboundRequest{
		Tag: i,
		Operation: serial.ToTypedMessage(&boundService.RemoveUserOperation{
			Email: u.Email,
		}),
	})
	return resp, err
}

// QueryStats get all stats
func (v *Client) QueryStats(c statsService.StatsServiceClient) (string, error) {
	r := &statsService.QueryStatsRequest{}
	rs := "pattern: \"\" reset: false"
	if err := proto.UnmarshalText(rs, r); err != nil {
		return "", err
	}
	resp, err := c.QueryStats(context.Background(), r)
	if err != nil {
		return "", err
	}
	return proto.MarshalTextString(resp), nil
}

// GetStats get state for user
func (v *Client) GetStats(c statsService.StatsServiceClient, u User) (string, error) {
	r := &statsService.GetStatsRequest{}
	rs := fmt.Sprintf("name: \"user>>>%v>>>traffic>>>downlink\" reset: false", u.Email)
	if err := proto.UnmarshalText(rs, r); err != nil {
		return "", err
	}
	resp, err := c.GetStats(context.Background(), r)
	if err != nil {
		return "", err
	}
	return proto.MarshalTextString(resp), nil
}

// ExampleV2rayClient example
func ExampleV2rayClient() {
	v := Client{
		Host: "127.0.0.1",
		Port: 10085,
	}
	u := User{
		Level:   0,
		Email:   "t@t.tt",
		UUID:    "a994b3c1-c7cc-4868-8072-c93e491bba0b",
		AlterID: 64,
	}
	c, err := v.GetConnection()
	if err != nil {
		panic(err)
	}
	hsClient := v.GetBondClient(c)
	v.AddUser(hsClient, "default", u)
	qsClient := v.GetStatClient(c)
	fmt.Println(v.QueryStats(qsClient))
	fmt.Println(v.GetStats(qsClient, u))
	v.RemoveUser(hsClient, "default", u)
	// Output: true
}
