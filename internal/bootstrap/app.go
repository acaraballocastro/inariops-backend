package bootstrap

import (
	"context"
	"database/sql"
)

type App struct {
	Repositories *Repositories
	Services     *Services
}

func New(
	db *sql.DB,
) *App {

	repositories := newRepositories(db)

	services := newServices(repositories)

	return &App{

		Repositories: repositories,

		Services: services,
	}
}

func (a *App) Start(
	ctx context.Context,
) {

	a.StartWorkers(ctx)
}
