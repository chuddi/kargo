package promotions

import (
	"context"
	"errors"

	kargoapi "github.com/akuity/kargo/api/v1alpha1"
	"github.com/akuity/kargo/pkg/database"
	"github.com/akuity/kargo/pkg/promotion"
)

// targetStore is the slice of database.Store the Promotion reconciler uses.
type targetStore interface {
	GetTarget(context.Context, string, string) (*kargoapi.Target, error)
}

// storeTargetGetter resolves Targets from the database in the shape
// promotion.ResolveTargetContext expects: a missing Target is nil, not an
// error, so that the reconciler can tell a Target that does not exist from a
// database it could not reach.
type storeTargetGetter struct {
	store targetStore
}

// NewTargetGetter returns a promotion.TargetGetter over the store, or nil
// when there is no store, which ResolveTargetContext treats as Targets being
// unavailable.
func NewTargetGetter(store database.Store) promotion.TargetGetter {
	if store == nil {
		return nil
	}
	return storeTargetGetter{store: store}
}

func (g storeTargetGetter) GetTarget(
	ctx context.Context,
	project string,
	name string,
) (*kargoapi.Target, error) {
	target, err := g.store.GetTarget(ctx, project, name)
	if errors.Is(err, database.ErrNotFound) {
		return nil, nil
	}
	return target, err
}
