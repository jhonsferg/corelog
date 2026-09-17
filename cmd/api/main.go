package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/jhonsferg/corelog/internal/platform/config"
	"github.com/jhonsferg/corelog/internal/platform/database"
	"github.com/jhonsferg/corelog/internal/platform/httpserver"
	"github.com/jhonsferg/corelog/internal/platform/logger"
	"github.com/jhonsferg/corelog/internal/team"
	teampostgres "github.com/jhonsferg/corelog/internal/team/adapter/out/postgres"
	teamsqlite "github.com/jhonsferg/corelog/internal/team/adapter/out/sqlite"
	teamport "github.com/jhonsferg/corelog/internal/team/application/port"
	"github.com/jhonsferg/corelog/internal/ticket"
	ticketdirectory "github.com/jhonsferg/corelog/internal/ticket/adapter/out/directory"
	ticketmemory "github.com/jhonsferg/corelog/internal/ticket/adapter/out/memory"
	ticketpostgres "github.com/jhonsferg/corelog/internal/ticket/adapter/out/postgres"
	ticketsqlite "github.com/jhonsferg/corelog/internal/ticket/adapter/out/sqlite"
	ticketport "github.com/jhonsferg/corelog/internal/ticket/application/port"
	"github.com/jhonsferg/corelog/internal/user"
	userpostgres "github.com/jhonsferg/corelog/internal/user/adapter/out/postgres"
	usersqlite "github.com/jhonsferg/corelog/internal/user/adapter/out/sqlite"
	userport "github.com/jhonsferg/corelog/internal/user/application/port"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal startup error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg.Server.Env)

	var db *sqlx.DB
	switch cfg.Database.Driver {
	case "sqlite":
		db, err = database.NewSQLite(cfg.Database)
	default:
		db, err = database.NewPostgres(cfg.Database)
	}
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("connected to database", "driver", cfg.Database.Driver)

	var userRepo userport.Repository
	var teamRepo teamport.Repository
	var ticketDBRepo ticketport.Repository
	switch cfg.Database.Driver {
	case "sqlite":
		userRepo = usersqlite.NewUserRepository(db)
		teamRepo = teamsqlite.NewTeamRepository(db)
		ticketDBRepo = ticketsqlite.NewTicketRepository(db)
	default:
		userRepo = userpostgres.NewUserRepository(db)
		teamRepo = teampostgres.NewTeamRepository(db)
		ticketDBRepo = ticketpostgres.NewTicketRepository(db)
	}

	userModule := user.NewModule(userRepo, cfg.Auth.JWTSecret, cfg.Auth.JWTExpiration)
	teamModule := team.NewModule(teamRepo, cfg.Auth.JWTSecret)

	ticketRepo := ticketDBRepo
	if cfg.Features.TicketRepository == "memory" {
		ticketRepo = ticketmemory.NewTicketRepository()
		log.Warn("ticket module using in-memory repository")
	}
	userDirectory := ticketdirectory.NewUserDirectory(db)
	ticketModule := ticket.NewModule(ticketRepo, userDirectory, cfg.Auth.JWTSecret)

	router := httpserver.NewRouter(cfg.Server)
	router.Route("/api/v1", func(r chi.Router) {
		httpserver.Mount(r, userModule.Handler, teamModule.Handler, ticketModule.Handler)
	})

	server := httpserver.New(cfg.Server, router)

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("server listening", "addr", server.Addr(), "env", cfg.Server.Env)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		log.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return err
		}
		log.Info("server stopped gracefully")
	}

	return nil
}
