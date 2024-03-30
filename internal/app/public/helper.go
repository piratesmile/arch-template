package public

import (
	"arch-template/ent"
	"arch-template/internal/app/entity"
	"arch-template/pkg/errdefs"
	"context"
)

func AuthIDFromContext(ctx context.Context) (uint, error) {
	uid := ctx.Value(CtxAuthIDKey)
	if v, ok := uid.(uint); ok && v > 0 {
		return v, nil
	}
	return 0, errdefs.ErrUnauthorized
}

func AuthModelFromContext(ctx context.Context) (*entity.User, error) {
	m := ctx.Value(CtxAuthModelKey)
	if v, ok := m.(*entity.User); ok && v.ID > 0 {
		return v, nil
	}
	return nil, errdefs.ErrUnauthorized
}

func ReplaceEntNotFoundError(err error, replace error) error {
	if ent.IsNotFound(err) {
		return replace
	}
	return err
}
