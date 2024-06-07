package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/datasource"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model/user"
)

type UserRepository struct {
	*datasource.DataSource
}

var (
	_ domain.UserRepository = (*UserRepository)(nil)
	_ tx.DataSource         = (*UserRepository)(nil)
)

func NewUserRepository(ds *datasource.DataSource) *UserRepository {
	return &UserRepository{DataSource: ds}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	return r.DataSource.TxOrPlain(ctx).User.Create().
		SetID(user.ID).
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetProfileURL(user.ProfileURL).
		SetRole(uint8(user.Role)).
		Exec(ctx)
}

func (r *UserRepository) FetchByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	entity, err := r.DataSource.TxOrPlain(ctx).User.Get(ctx, id)
	if err != nil {
		if model.IsNotFound(err) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	user := domain.User{
		ID:         entity.ID,
		Username:   entity.Username,
		Email:      entity.Email,
		ProfileURL: entity.ProfileURL,
		Role:       domain.UserRole(entity.Role),
	}

	return user, nil
}

func (r *UserRepository) FetchByUsername(ctx context.Context, username string) (domain.User, error) {
	entity, err := r.DataSource.TxOrPlain(ctx).User.
		Query().
		Where(user.Username(username)).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	user := domain.User{
		ID:         entity.ID,
		Username:   entity.Username,
		Email:      entity.Email,
		ProfileURL: entity.ProfileURL,
		Role:       domain.UserRole(entity.Role),
	}

	return user, nil
}

func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return r.DataSource.TxOrPlain(ctx).User.
		Query().
		Where(user.Username(username)).
		Exist(ctx)
}
