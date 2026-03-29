package app

import (
	"github.com/DangeL187/erax"

	"back/internal/features/subs/handler"
	"back/internal/features/subs/infra"
	"back/internal/features/subs/usecase"
	"back/internal/infra/database"
	"back/internal/shared/config"
)

type App struct {
	Config      *config.Config
	SubsHandler *handler.SubsHandler
}

func NewApp() (*App, error) {
	app := &App{}

	var err error
	app.Config, err = config.NewConfig()
	if err != nil {
		return nil, erax.Wrap(err, "failed to load config")
	}

	db, err := database.NewPostgres(app.Config)
	if err != nil {
		return nil, erax.Wrap(err, "failed to connect to DB")
	}

	crudlRepo := infra.NewCrudlRepo(db)
	sumRepo := infra.NewSumRepo(db)

	crudlUseCase := usecase.NewCrudlUseCase(crudlRepo)
	sumUseCase := usecase.NewSumUseCase(sumRepo)
	app.SubsHandler = handler.NewSubsHandler(crudlUseCase, sumUseCase)

	return app, nil
}
