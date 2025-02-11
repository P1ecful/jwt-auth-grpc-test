package tests

import (
	"context"
	gen "github.com/P1ecful/pkg/gen/grpc/auth"
	authmock "github.com/P1ecful/pkg/gen/grpc/auth/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestLogin_Happy(t *testing.T) {
	type mockBehavior func(
		mockClient *authmock.MockAuthClient,
		req *gen.LoginRequest,
		response error,
	)

	cases := []struct {
		name             string
		request          *gen.LoginRequest
		expectedResponse error
		mockBehavior     mockBehavior
	}{
		{
			name: "Successful",
			request: &gen.LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedResponse: nil,
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.LoginRequest,
				response error,
			) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(nil, response)
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := authmock.NewMockAuthClient(ctrl)
			cs.mockBehavior(client, cs.request, cs.expectedResponse)
			_, actual := client.Login(context.Background(), cs.request)

			assert.Equal(t, cs.expectedResponse, actual)
		})
	}
}

func TestLogin_Fail(t *testing.T) {
	type mockBehavior func(
		mockClient *authmock.MockAuthClient,
		req *gen.LoginRequest,
		response error,
	)

	cases := []struct {
		name             string
		request          *gen.LoginRequest
		expectedResponse error
		mockBehavior     mockBehavior
	}{
		{
			name: "Login without password",
			request: &gen.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			expectedResponse: status.Error(codes.InvalidArgument, "password is required"),
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.LoginRequest,
				response error,
			) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(nil, response)
			},
		},
		{
			name: "Login without email",
			request: &gen.LoginRequest{
				Email:    "",
				Password: "password",
			},
			expectedResponse: status.Error(codes.InvalidArgument, "email is required"),
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.LoginRequest,
				response error,
			) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(nil, response)
			},
		},
		{
			name: "Login without data",
			request: &gen.LoginRequest{
				Email:    "",
				Password: "",
			},
			expectedResponse: status.Error(codes.InvalidArgument, "email is required"),
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.LoginRequest,
				response error,
			) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(nil, response)
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := authmock.NewMockAuthClient(ctrl)
			cs.mockBehavior(client, cs.request, cs.expectedResponse)
			_, actual := client.Login(context.Background(), cs.request)

			assert.Equal(t, cs.expectedResponse, actual)
		})
	}
}
