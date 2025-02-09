package pgx

import (
	"context"
	"errors"
	"github.com/P1ecful/jwt-grpc-test/internal/model/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrEmailExists  = errors.New("email is already exists")
	ErrUserNotFound = errors.New("user not found")
)

const (
	createUserQuery      = `INSERT INTO users (email, password) VALUES (@email, @password)`
	checkUserExistsQuery = `SELECT * FROM users WHERE email = @email `
)

type PGX struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPGX(logger *zap.Logger, path string) *PGX {
	pool, err := pgxpool.New(context.Background(), path)

	if err != nil {
		logger.Fatal("unable to create connection pool", zap.Error(err))
	}

	logger.Info("database initialized and successfully connected")
	return &PGX{
		pool:   pool,
		logger: logger,
	}
}

func (p *PGX) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *PGX) Disconnect() {
	p.pool.Close()
}

func (p *PGX) CreateNewUser(ctx context.Context, user dto.User) error {
	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
	}

	_, err := p.pool.Exec(ctx, createUserQuery, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrEmailExists
			}
		}

		return err
	}

	return nil
}

func (p *PGX) User(ctx context.Context, email string) (dto.User, error) {
	var usr dto.User

	if err := p.pool.QueryRow(ctx, checkUserExistsQuery,
		pgx.NamedArgs{"email": email}).Scan(&usr.Email, &usr.Password); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			p.logger.Debug("user not found", zap.String("email", email))

			return dto.User{}, ErrUserNotFound
		}

		return dto.User{}, err
	}

	return usr, nil
}
