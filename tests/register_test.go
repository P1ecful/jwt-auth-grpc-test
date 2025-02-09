package tests

import (
	gen "github.com/P1ecful/pkg/gen/grpc/auth"
	authmock "github.com/P1ecful/pkg/gen/grpc/auth/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRegister_Handler(t *testing.T) {
	type mockBehavior func(
		mockClient *authmock.MockAuthClient,
		req *gen.RegisterRequest,
		expectedResponse *gen.RegisterResponse,
	)

	//	name:     "Register without password",
	//	email:    "TestMail@gmail.com",
	//	password: "",
	//	except:   "password is required",
	//},
	//{
	//	name:     "Register without email",
	//	email:    "",
	//	password: "TestPassword",
	//	except:   "email is required",
	//},
	//{
	//	name:     "Register without data",
	//	email:    "",
	//	password: "",
	//	except:   "email is required",
	//},

	// !TODO more mock_test cases
	cases := []struct {
		name             string
		requestBody      map[string]string
		exceptedError    error
		expectedRequest  *gen.RegisterRequest
		expectedResponse *gen.RegisterResponse
		mockBehavior     mockBehavior
	}{
		{
			name: "Register Happy Case",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password",
			},
			exceptedError: nil,
			expectedRequest: &gen.RegisterRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedResponse: &gen.RegisterResponse{
				Status: "Successful",
			},
			mockBehavior: func(
				mockClient *authmock.MockAuthClient,
				req *gen.RegisterRequest,
				expectedResponse *gen.RegisterResponse,
			) {
				mockClient.EXPECT().Register(gomock.Any(), req).Return(expectedResponse, nil)
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := authmock.NewMockAuthClient(ctrl)
			cs.mockBehavior(client, cs.expectedRequest, cs.expectedResponse)

			assert.Equal(t, cs.expectedResponse, cs.requestBody)
		})
	}
}

// !TODO login and get_data
//func TestLogin(t *testing.T) {
//	type mockBehavior func(
//		mockClient *authmock.MockAuthClient,
//		req *gen.LoginRequest,
//		expectedResponse *gen.LoginResponse,
//	)
//
//	cases := []struct {
//		name             string
//		requestBody      map[string]string
//		exceptedError    error
//		expectedRequest  *gen.LoginRequest
//		expectedResponse *gen.LoginResponse
//		mockBehavior     mockBehavior
//	}{}
//}
