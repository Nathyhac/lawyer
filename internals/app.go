package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/Nathac/go-api/internals/handlers"
	"github.com/Nathac/go-api/internals/store"
	"github.com/Nathac/go-api/migration"
)

type Application struct {
	LawyerHandler *handlers.LawyerHandler
	UserHandler   *handlers.UserHandler
	TokenHandler  *handlers.TokenHandler
	Logger        *log.Logger
	DB            *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "INFO", log.Ldate|log.Ltime)
	Pgdb, err := store.Open()
	if err != nil {
		log.Fatalf("error happend while opening a connection to db %v", err)
	}

	//---STORES---//
	LawyerStore := store.NewPostgresDB(Pgdb)
	UserStore := store.NewUserpostgresDB(Pgdb)
	TokenStore := store.NewPostgresTokenStore(Pgdb)

	//--HANDLERS--//
	lawyerHandler := handlers.NewLawyerHandler(LawyerStore, logger)
	userHandler := handlers.NewUserHandler(UserStore)
	tokenHandler := handlers.NewTokenHandler(UserStore, TokenStore)

	err = store.MigrateFs(Pgdb, migration.FS, ".")
	if err != nil {
		panic(err)
	}

	app := &Application{
		Logger:        logger,
		DB:            Pgdb,
		LawyerHandler: lawyerHandler,
		UserHandler:   userHandler,
		TokenHandler:  tokenHandler,
	}
	return app, nil
}
