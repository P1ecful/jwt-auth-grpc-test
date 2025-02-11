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

func TestGetData_Fail(t *testing.T) {
	type mockBehavior func(
		mockClient *authmock.MockAuthClient,
		req *gen.GetDataFromAccessTokenRequest,
		response error,
	)

	cases := []struct {
		name             string
		request          *gen.GetDataFromAccessTokenRequest
		expectedResponse error
		mockBehavior     mockBehavior
	}{
		{
			name: "Without token",
			request: &gen.GetDataFromAccessTokenRequest{
				Token: "",
			},
			expectedResponse: status.Error(codes.InvalidArgument, "token is required"),
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.GetDataFromAccessTokenRequest,
				response error,
			) {
				mockClient.EXPECT().GetDataFromAccessToken(gomock.Any(), req).Return(nil, response)
			},
		},
		{
			name: "Without token",
			request: &gen.GetDataFromAccessTokenRequest{
				Token: "",
			},
			expectedResponse: status.Error(codes.Unauthenticated, "error getting data from access token"),
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.GetDataFromAccessTokenRequest,
				response error,
			) {
				mockClient.EXPECT().GetDataFromAccessToken(gomock.Any(), req).Return(nil, response)
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := authmock.NewMockAuthClient(ctrl)
			cs.mockBehavior(client, cs.request, cs.expectedResponse)
			_, actual := client.GetDataFromAccessToken(context.Background(), cs.request)

			assert.Equal(t, cs.expectedResponse, actual)
		})
	}
}
