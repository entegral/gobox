package dynamo

import (
	"context"

	ttypes "github.com/entegral/gobox/types"
)

func (m *DiLink[T0, T1]) LoadEntity0s(ctx context.Context, linkWrapper ttypes.Linkable) ([]T0, error) {
	loaded, err := m.LoadEntity1(ctx)
	if err != nil {
		return nil, err
	}
	if !loaded {
		return nil, ErrEntityNotFound[T1]{Entity: m.Entity1}
	}
	links, err := FindLinksByEntity1[T1, *DiLink[T0, T1]](ctx, m.Entity1, linkWrapper.Type())
	if err != nil {
		return nil, err
	}
	var entities []T0
	for _, link := range links {
		loaded, err := link.LoadEntity0(ctx)
		if err != nil {
			return nil, err
		}
		if loaded {
			entities = append(entities, link.Entity0)
		}
	}
	return entities, nil
}
