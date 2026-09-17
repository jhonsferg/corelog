package user

import (
	"time"

	httpadapter "github.com/jhonsferg/corelog/internal/user/adapter/in/http"
	"github.com/jhonsferg/corelog/internal/user/application/port"
	"github.com/jhonsferg/corelog/internal/user/application/usecase"
)

type Module struct {
	Service port.Service
	Handler *httpadapter.Handler
}

func NewModule(repo port.Repository, jwtSecret string, jwtExpiration time.Duration) *Module {
	service := usecase.NewUserService(repo)
	handler := httpadapter.NewHandler(service, jwtSecret, jwtExpiration)

	return &Module{
		Service: service,
		Handler: handler,
	}
}
