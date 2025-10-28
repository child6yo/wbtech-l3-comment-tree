package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/child6yo/wbtech-l3-comment-tree/internal/controller"
	"github.com/child6yo/wbtech-l3-comment-tree/internal/logger"
	"github.com/child6yo/wbtech-l3-comment-tree/internal/repository/postgres"
	"github.com/child6yo/wbtech-l3-comment-tree/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

const (
	addCommentRoute      = "/comments"
	getCommentsTreeRoute = "/comments"
	deleteComment        = "/comments/:id"
)

type appConfig struct {
	address string

	pgHost     string
	pgPort     string
	pgUsername string
	pgDBName   string
	pgPassword string
	pgSSLMode  string
}

func initConfig(envFilePath, envPrefix string) (*appConfig, error) {
	appConfig := &appConfig{}

	cfg := config.New()

	err := cfg.LoadEnvFiles(envFilePath)
	cfg.EnableEnv(envPrefix)

	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	appConfig.address = cfg.GetString("ADDR")

	// PostgreSQL
	appConfig.pgHost = cfg.GetString("PG_HOST")
	appConfig.pgPort = cfg.GetString("PG_PORT")
	appConfig.pgPassword = cfg.GetString("PG_PASSWORD")
	appConfig.pgUsername = cfg.GetString("PG_USER")
	appConfig.pgDBName = cfg.GetString("PG_DB")
	appConfig.pgSSLMode = cfg.GetString("PG_SSLMODE")

	return appConfig, nil
}

func main() {
	var wg sync.WaitGroup

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	zlog.InitConsole()
	lgr := zlog.Logger

	cfg, err := initConfig(".env", "")
	if err != nil {
		lgr.Fatal().Err(err).Send()
	}

	db, err := postgres.NewMSPostgresDB(
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.pgHost, cfg.pgPort, cfg.pgUsername, cfg.pgDBName, cfg.pgPassword, cfg.pgSSLMode),
	)
	if err != nil {
		lgr.Fatal().Err(err).Send()
	}

	cr := postgres.NewCommentsRepository(db)

	cs := usecase.NewCommentsService(cr)

	cc := controller.NewCommentsController(cs)
	mdlw := controller.NewMiddleware(logger.NewLoggerAdapter(lgr))

	srv := ginext.New("")
	srv.Use(ginext.Logger(), ginext.Recovery(), mdlw.ErrHandlingMiddleware(), cors.Default())
	srv.POST(addCommentRoute, cc.NewComment)
	srv.GET(getCommentsTreeRoute, cc.GetComments)
	srv.DELETE(deleteComment, cc.DeleteComments)

	httpServer := &http.Server{
		Addr:    cfg.address,
		Handler: srv,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lgr.Err(err).Send()
		}
	}()

	<-ctx.Done()
	lgr.Info().Msg("shutting down gracefully...")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		lgr.Err(err).Send()
	}

	if err := db.Master.Close(); err != nil {
		lgr.Err(err).Send()
	}

	wg.Wait()

	lgr.Info().Msg("app exited")
}
