package api

import (
	"awesomeProject/server/api/proto/generated"
	"awesomeProject/server/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server
type Server struct {
	api.UnimplementedAuthenticationServer
}

func (s *Server) SignUp(ctx context.Context, login *api.RequestLogin) (*api.ResponseLogin, error) {
	if login.UserName == "" || login.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Некоректные данные")
	}

	sessionKey, err := service.SignUp(login)
	if err != nil {
		return nil, err
	}
	return &api.ResponseLogin{SessionKey: sessionKey}, nil
}

func (s *Server) SignIn(ctx context.Context, login *api.RequestLogin) (*api.ResponseLogin, error) {
	if login.SessionKey == "" {
		if login.UserName == "" || login.Password == "" {
			return nil, status.Errorf(codes.InvalidArgument, "Некоректные данные")
		}
	}

	sessionKey, err := service.SignIn(login)
	if err != nil {
		return nil, err
	}
	return &api.ResponseLogin{SessionKey: sessionKey}, nil
}

func (s *Server) mustEmbedUnimplementedAuthenticationServer() {
	panic("implement me")
}
