package team

import (
	httpadapter "github.com/jhonsferg/corelog/internal/team/adapter/in/http"
	"github.com/jhonsferg/corelog/internal/team/application/port"
	"github.com/jhonsferg/corelog/internal/team/application/usecase"
)

type Module struct {
	Service port.Service
	Handler *httpadapter.Handler
}

func NewModule(repo port.Repository, jwtSecret string) *Module {
	service := usecase.NewTeamService(repo)
	handler := httpadapter.NewHandler(service, jwtSecret)

	return &Module{
		Service: service,
		Handler: handler,
	}
}
