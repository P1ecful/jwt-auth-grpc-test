package storage

import (
	"context"
	"github.com/P1ecful/jwt-grpc-test/internal/model/dto"
)

type PostgresStorage interface {
	Ping(ctx context.Context) error
	CreateNewUser(ctx context.Context, user dto.User) error
	User(ctx context.Context, email string) (dto.User, error)
	Disconnect()
}
