package bootstrap

import (
	"database/sql"
)

type App struct {
	Repositories *Repositories
	Services     *Services
	Handlers     *Handlers
}

func New(db *sql.DB) *App {

	repositories := newRepositories(db)

	services := newServices(repositories)

	handlers := newHandlers(services)

	return &App{
		Repositories: repositories,
		Services:     services,
		Handlers:     handlers,
	}
}
