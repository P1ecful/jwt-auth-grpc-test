package tests

import (
	gen "github.com/P1ecful/pkg/gen/grpc/auth"
	authmock "github.com/P1ecful/pkg/gen/grpc/auth/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestRegister_Happy(t *testing.T) {
	type mockBehavior func(
		mockClient *authmock.MockAuthClient,
		request *gen.RegisterRequest,
		response *gen.RegisterResponse,
	)

	cases := []struct {
		name             string
		Request          *gen.RegisterRequest
		exceptedResponse *gen.RegisterResponse
		mockBehavior     mockBehavior
	}{
		{
			name: "Register Happy Case",
			Request: &gen.RegisterRequest{
				Email:    "test@example.com",
				Password: "password",
			},

			exceptedResponse: &gen.RegisterResponse{
				Status: "Successful",
			},
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.RegisterRequest,
				response *gen.RegisterResponse,
			) {
				mockClient.EXPECT().Register(gomock.Any(), req).Return(response, nil)
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := authmock.NewMockAuthClient(ctrl)
			cs.mockBehavior(client, cs.Request, cs.exceptedResponse)

			actual, _ := client.Register(context.Background(), cs.Request)
			assert.Equal(t, cs.exceptedResponse, actual)
		})
	}
}

func TestRegister_Fail(t *testing.T) {
	type mockBehavior func(
		mockClient *authmock.MockAuthClient,
		request *gen.RegisterRequest,
		response error,
	)

	cases := []struct {
		name             string
		Request          *gen.RegisterRequest
		exceptedResponse error
		mockBehavior     mockBehavior
	}{
		{
			name: "Register without password",
			Request: &gen.RegisterRequest{
				Email:    "test@example.com",
				Password: "",
			},
			exceptedResponse: nil, //status.Error(codes.InvalidArgument, "password is required"),
			mockBehavior: func(mockClient *authmock.MockAuthClient,
				req *gen.RegisterRequest,
				response error,
			) {
				mockClient.EXPECT().Register(gomock.Any(), req).Return(nil, response)
			},
		},
		{
			name: "Register without email",

			Request: &gen.RegisterRequest{
				Email:    "",
				Password: "password",
			},

			exceptedResponse: status.Error(codes.InvalidArgument, "email is required"),
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.RegisterRequest,
				response error,
			) {
				mockClient.EXPECT().Register(gomock.Any(), req).Return(nil, response)
			},
		},
		{
			name: "Register without data",
			Request: &gen.RegisterRequest{
				Email:    "",
				Password: "",
			},
			exceptedResponse: status.Error(codes.InvalidArgument, "password is required"),
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.RegisterRequest,
				response error,
			) {
				mockClient.EXPECT().Register(gomock.Any(), req).Return(nil, response)
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := authmock.NewMockAuthClient(ctrl)
			cs.mockBehavior(client, cs.Request, cs.exceptedResponse)
			_, actual := client.Register(context.Background(), cs.Request)

			assert.Equal(t, cs.exceptedResponse, actual)
		})
	}
}
