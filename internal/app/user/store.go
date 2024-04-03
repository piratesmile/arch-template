package user

import (
	"arch-template/ent/predicate"
	"arch-template/ent/user"
	"arch-template/internal/app/entity"
	"arch-template/internal/app/public"
	"arch-template/pkg/errdefs"
	"context"
)

type Store struct {
	*public.EntStore
}

func NewRepository(baseStore *public.EntStore) *Store {
	return &Store{EntStore: baseStore}
}

func (u *Store) FetchByID(ctx context.Context, id uint) (*entity.User, error) {
	return u.only(ctx, user.ID(id))
}

func (u *Store) FetchByUserName(ctx context.Context, username string) (*entity.User, error) {
	return u.only(ctx, user.Username(username))
}

func (u *Store) Create(ctx context.Context, user *entity.User) error {
	userM, err := u.Client(ctx).User.
		Create().
		SetUsername(user.UserName).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		return err
	}
	(*user).ID = userM.ID
	return nil
}

func (u *Store) only(ctx context.Context, conditions ...predicate.User) (*entity.User, error) {
	entUser, err := u.Client(ctx).User.Query().Where(conditions...).Only(ctx)
	if err != nil {
		return nil, public.ReplaceEntNotFoundError(err, errdefs.ErrResourceNotFound)
	}
	return entity.UserFromEntModel(entUser), nil
}

func (u *Store) all(ctx context.Context, conditions ...predicate.User) ([]*entity.User, error) {
	list, err := u.Client(ctx).User.Query().Where(conditions...).All(ctx)
	if err != nil {
		return nil, err
	}
	var res = make([]*entity.User, len(list))
	for i := 0; i < len(list); i++ {
		res[i] = entity.UserFromEntModel(list[i])
	}
	return res, nil
}
