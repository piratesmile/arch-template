package public

import (
	"arch-template/ent"
	"context"
	"fmt"
)

type EntRepository struct {
	entClient *ent.Client
}

func NewBaseRepository(client *ent.Client) *EntRepository {
	return &EntRepository{entClient: client}
}

func (b *EntRepository) Client(ctx context.Context) *ent.Client {
	if tx := ent.TxFromContext(ctx); tx != nil {
		return tx.Client()
	}
	return b.entClient
}

func (b *EntRepository) Tx(ctx context.Context, fn func(context.Context) error) error {
	if ent.TxFromContext(ctx) != nil {
		return fn(ctx)
	}
	tx, err := b.entClient.Tx(ctx)
	if err != nil {
		return err
	}
	txContext := ent.NewTxContext(ctx, tx)
	if err := fn(txContext); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rollback transaction err: %v", err, rerr)
		}
		return err
	}
	return tx.Commit()
}
