package auth

import (
	"context"
	"github.com/P1ecful/jwt-grpc-test/internal/service"
	gen "github.com/P1ecful/pkg/gen/grpc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	gen.UnimplementedAuthServer
	service service.Service
	logger  *zap.Logger
}

func NewGRPCServer(logger *zap.Logger, service service.Service) *GRPCServer {
	return &GRPCServer{
		service: service,
		logger:  logger,
	}
}

func (gs *GRPCServer) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	if req.Email == "" {
		gs.logger.Debug("Empty email")
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		gs.logger.Debug("Password is required")
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if err := gs.service.Register(ctx, req.GetEmail(), req.GetPassword()); err != nil {
		gs.logger.Debug("Register error", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	gs.logger.Debug("Register success", zap.Any("email", req.Email))
	return &gen.RegisterResponse{Status: "Successful"}, nil
}

func (gs *GRPCServer) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	if req.Email == "" {
		gs.logger.Debug("Empty email")
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		gs.logger.Debug("Password is required")
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := gs.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		gs.logger.Debug("Login error", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	gs.logger.Debug("Login success", zap.Any("Email", req.Email))
	return &gen.LoginResponse{Token: token}, nil
}

func (gs *GRPCServer) GetDataFromAccessToken(_ context.Context, req *gen.GetDataFromAccessTokenRequest) (*gen.GetDataFromAccessTokenResponse, error) {
	if req.Token == "" {
		gs.logger.Debug("Empty token")
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	email, err := gs.service.GetDataFromAccessToken(req.GetToken())
	if err != nil {
		gs.logger.Debug("Get data error", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "error getting data from access token")
	}

	return &gen.GetDataFromAccessTokenResponse{Email: email}, nil
}
