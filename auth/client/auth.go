package client

import (
	"context"
	"log"
	"time"

	pb "github.com/dittrichlucas/poc-grpc-bookstore/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthClient struct {
	service pb.AuthServiceClient
	token   string
	// username string
	// password string
}

type AuthClientInterceptor struct { //AuthInterceptor
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthClient(cc *grpc.ClientConn, token string) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{service, token}
}

func NewAuthClientInterceptor(
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) (*AuthClientInterceptor, error) {
	interceptor := &AuthClientInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	if err := interceptor.scheduleRefreshToken(refreshDuration); err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (i *AuthClientInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> unary interceptor: %s", method)

		if i.authMethods[method] {
			return invoker(i.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *AuthClientInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		if i.authMethods[method] {
			return streamer(i.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (i *AuthClientInterceptor) attachToken(ctx context.Context) context.Context {
	log.Println("--> token: " + i.accessToken)
	return metadata.AppendToOutgoingContext(ctx, "authorization", i.accessToken)
}

// func (c *AuthClient) Login() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	req := &pb.LoginRequest{
// 		Username: c.username,
// 		Password: c.password,
// 	}

// 	res, err := c.service.Login(ctx, req)
// 	if err != nil {
// 		return "", err
// 	}

// 	return res.GetAccessToken(), nil
// }

func (c *AuthClient) Auth() string {
	return c.token
}

func (i *AuthClientInterceptor) refreshToken() error {
	// accessToken, err := i.authClient.Login()
	// if err != nil {
	// 	return err
	// }

	i.accessToken = i.authClient.Auth() // "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzI1MjQzNjMsInVzZXJuYW1lIjoidXNlcjEiLCJyb2xlIjoidXNlciJ9.iGils6zqSPkdkWh2_GOeCLqX2QLYDCoh3X6jkFHSq6g"
	log.Printf("token refreshed: %v", i.accessToken)

	return nil
}

func (i *AuthClientInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	if err := i.refreshToken(); err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			if err := i.refreshToken(); err != nil {
				log.Println("--> teste goroutine")
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}
