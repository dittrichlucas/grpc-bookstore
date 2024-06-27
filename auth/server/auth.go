package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dittrichlucas/poc-grpc-bookstore/proto"
	s "github.com/dittrichlucas/poc-grpc-bookstore/service"
)

type AuthService struct {
	userStore  s.UserStore
	jwtManager s.JWTManager
}

func NewAuthService(userStore s.UserStore, jwtManager s.JWTManager) *AuthService {
	return &AuthService{userStore, jwtManager}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userStore.Find(req.GetUsername())
	// fmt.Printf("user: %v", req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.LoginResponse{AccessToken: token}
	return res, nil
}
