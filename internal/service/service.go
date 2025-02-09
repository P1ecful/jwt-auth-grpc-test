package service

import (
	"context"
	"github.com/P1ecful/jwt-grpc-test/internal/config"
	"github.com/P1ecful/jwt-grpc-test/internal/model/dto"
	"github.com/P1ecful/jwt-grpc-test/internal/service/jwt"
	"github.com/P1ecful/jwt-grpc-test/internal/storage"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (string, error)
	GetDataFromAccessToken(token string) (string, error)
}

type Auth struct {
	logger  *zap.Logger
	Storage storage.PostgresStorage
	config  config.ServiceConfig
}

func NewAuth(logger *zap.Logger,
	postgres storage.PostgresStorage,
	config config.ServiceConfig) *Auth {
	return &Auth{
		logger:  logger,
		Storage: postgres,
		config:  config,
	}
}

func (s *Auth) Register(ctx context.Context, email string, password string) error {
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err := s.Storage.CreateNewUser(ctx, dto.User{Email: email, Password: hashPwd}); err != nil {
		s.logger.Debug("Failed to register user",
			zap.String("email", email),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (s *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.Storage.User(ctx, email)
	if err != nil {
		s.logger.Debug("Failed to get user",
			zap.String("email", email),
			zap.Error(err),
		)

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		s.logger.Debug("Invalid password",
			zap.String("email", email),
			zap.Error(err),
		)

		return "", err
	}

	token, err := jwt.GenerateTokens(email, s.config.SecretKey, s.config.AccessTokenTTL)
	if err != nil {
		s.logger.Debug("Failed to generate token", zap.Error(err))
		return "", err
	}

	return token, nil
}

func (s *Auth) GetDataFromAccessToken(token string) (string, error) {
	claims := jwtlib.MapClaims{}
	_, err := jwtlib.ParseWithClaims(token, claims, func(token *jwtlib.Token) (interface{}, error) {
		return nil, nil
	})

	if err != nil {
		s.logger.Debug("Failed to parse token", zap.Error(err))
		return "", err
	}

	s.logger.Debug("Got data", zap.Any("claims", claims))
	return claims["email"].(string), nil
}
