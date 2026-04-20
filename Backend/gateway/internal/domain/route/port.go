package route

import "context"

type Repository interface {
	Save(ctx context.Context, route *Route) error

	FindByID(ctx context.Context, id RouteID) (*Route, error)

	FindByPath(ctx context.Context, path string, method string) (*Route, error)

	Delete(ctx context.Context, id RouteID) error
}
