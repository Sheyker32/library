package run

import (
	"context"
	"fmt"
	"library/internal/facade"
	"library/internal/handler"
	"library/internal/repository"
	"library/internal/usecase"
	"library/responder"
	"library/router"
	"library/server"

	"go.uber.org/zap"

	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/jmoiron/sqlx"
	"github.com/ptflp/godecoder"

	"golang.org/x/sync/errgroup"
)

const (
	NoError = iota
	InternalError
	GeneralError
)

// Application - интерфейс приложения
type Application interface {
	Runner
	Bootstraper
}

// Runner - интерфейс запуска приложения
type Runner interface {
	Run() int
}

// Bootstraper - интерфейс инициализации приложения
type Bootstraper interface {
	Bootstrap(options ...interface{}) Runner
}

// App - структура приложения
type App struct {
	logger *zap.Logger
	db     *sqlx.DB
	srv    *server.Server
	Sig    chan os.Signal
}

// NewApp - конструктор приложения
func NewApp(db *sqlx.DB, logger *zap.Logger) *App {
	return &App{db: db, logger: logger, Sig: make(chan os.Signal, 1)}
}

// Run - запуск приложения
func (a *App) Run() int {
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sigInt))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("app: server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return GeneralError
	}

	return NoError
}

func (a *App) Bootstrap(options ...interface{}) Runner {
	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})
	respond := responder.NewResponder(decoder, a.logger)

	authorRepo := repository.NewAuthorRepository(a.db)
	bookRepo := repository.NewBookRepository(a.db, authorRepo)
	userRepo := repository.NewUserRepository(a.db)
	rentRepo := repository.NewRentalRepository(a.db)

	userUC := usecase.NewUserUseCase(userRepo)
	authorUC := usecase.NewAuthorUseCase(authorRepo)
	bookUC := usecase.NewBookUseCase(bookRepo)
	rentUC := usecase.NewRentUseCase(rentRepo)

	facade := facade.NewLibraryFacade(a.db, authorUC, bookUC, rentUC, userUC)

	ctx := context.Background()
	err := facade.InitializeDataIfEmpty(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	authorHandler := handler.NewAuthorHandler(authorUC, respond)
	bookHandler := handler.NewBookHandler(bookUC, respond)
	userHandler := handler.NewUserHandler(userUC, respond)
	rentHandler := handler.NewRentHandler(facade, respond)

	r := router.NewApiRouter(authorHandler, bookHandler, rentHandler, userHandler)
	a.srv = server.NewServer(r)

	return a
}
