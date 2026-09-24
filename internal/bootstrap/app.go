package bootstrap

import (
	"database/sql"

	"inariops/internal/modules/auth"
)

type App struct {
	Repositories *Repositories
	Services     *Services
	Handlers     *Handlers
}

func New(
	db *sql.DB,
	jwtManager *auth.JWTManager,
) *App {

	repositories := newRepositories(db)

	services := newServices(
		repositories,
		jwtManager,
	)

	handlers := newHandlers(services)

	return &App{
		Repositories: repositories,
		Services:     services,
		Handlers:     handlers,
	}
}
