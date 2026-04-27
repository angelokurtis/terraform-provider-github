package datasource

import (
	"context"

	"github.com/google/go-github/v84/github"
	"github.com/sourcegraph/conc/pool"

	"github.com/angelokurtis/terraform-provider-github/internal/errors"
)

func fetchAll[T any, Slice ~[]T](
	ctx context.Context,
	fetch func(ctx context.Context, page int) (Slice, *github.Response, error),
) (Slice, error) {
	first, resp, err := fetch(ctx, 1)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make(Slice, 0, len(first))
	result = append(result, first...)

	if resp == nil || resp.LastPage <= 1 {
		return result, nil
	}

	p := pool.NewWithResults[Slice]().WithErrors().WithMaxGoroutines(50)

	for page := 2; page <= resp.LastPage; page++ {
		page := page

		p.Go(func() (Slice, error) {
			items, _, err := fetch(ctx, page)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return items, nil
		})
	}

	batches, err := p.Wait()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, batch := range batches {
		result = append(result, batch...)
	}

	return result, nil
}
