package test

import (
	"time"

	"github.com/DangeL187/erax"
	"gorm.io/gorm"

	"back/internal/features/subs/domain"
	"back/internal/features/subs/infra"
	"back/internal/infra/database"
	"back/internal/shared/config"
)

type DB struct {
	DB        *gorm.DB
	CrudlRepo domain.CrudlRepo
	SumRepo   domain.SumRepo
	Teardown  func()
}

func SetupTestDB() (*DB, error) {
	cfg := &config.Config{
		DBConnectTimeout: 1 * time.Minute,
		PostgresDSN:      "host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=disable",
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, erax.Wrap(err, "failed to connect to database")
	}

	cr := infra.NewCrudlRepo(db)
	sr := infra.NewSumRepo(db)

	return &DB{
		DB:        db,
		CrudlRepo: cr,
		SumRepo:   sr,
	}, nil
}
