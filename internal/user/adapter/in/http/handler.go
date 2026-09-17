package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
	"github.com/jhonsferg/corelog/internal/platform/authtoken"
	"github.com/jhonsferg/corelog/internal/platform/httpresponse"
	appmw "github.com/jhonsferg/corelog/internal/platform/middleware"
	"github.com/jhonsferg/corelog/internal/platform/validation"
	"github.com/jhonsferg/corelog/internal/user/application/port"
)

type Handler struct {
	service       port.Service
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewHandler(service port.Service, jwtSecret string, jwtExpiration time.Duration) *Handler {
	return &Handler{service: service, jwtSecret: jwtSecret, jwtExpiration: jwtExpiration}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	user, err := h.service.Register(r.Context(), port.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	user, err := h.service.Authenticate(r.Context(), port.AuthenticateInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	token, err := authtoken.Generate(h.jwtSecret, h.jwtExpiration, user.ID, string(user.Role))
	if err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindInternal, "failed to issue token", err))
		return
	}

	httpresponse.JSON(w, http.StatusOK, authResponse{Token: token, User: toUserResponse(user)})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmw.UserIDFromContext(r.Context())
	if !ok {
		httpresponse.Error(w, apperror.New(apperror.KindUnauthorized, "not authenticated"))
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toUserResponse(user))
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmw.UserIDFromContext(r.Context())
	if !ok {
		httpresponse.Error(w, apperror.New(apperror.KindUnauthorized, "not authenticated"))
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	user, err := h.service.UpdateProfile(r.Context(), userID, port.UpdateProfileInput{Name: req.Name, Email: req.Email})
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toUserResponse(user))
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmw.UserIDFromContext(r.Context())
	if !ok {
		httpresponse.Error(w, apperror.New(apperror.KindUnauthorized, "not authenticated"))
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	if err := h.service.ChangePassword(r.Context(), userID, port.ChangePasswordInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}); err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusNoContent, nil)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmw.UserIDFromContext(r.Context())
	if !ok {
		httpresponse.Error(w, apperror.New(apperror.KindUnauthorized, "not authenticated"))
		return
	}

	query := r.URL.Query().Get("q")

	users, err := h.service.SearchTeammates(r.Context(), userID, query)
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toUserResponses(users))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := port.SearchFilter{Query: q.Get("q")}
	if teamIDStr := q.Get("team_id"); teamIDStr != "" {
		teamID, err := uuid.Parse(teamIDStr)
		if err != nil {
			httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid team id"))
			return
		}
		filter.TeamID = &teamID
	}

	users, err := h.service.Search(r.Context(), filter)
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toUserResponses(users))
}

func (h *Handler) SetTeam(w http.ResponseWriter, r *http.Request) {
	targetID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid user id"))
		return
	}

	var req setTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	var teamID *uuid.UUID
	if req.TeamID != nil {
		parsed, err := uuid.Parse(*req.TeamID)
		if err != nil {
			httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid team id"))
			return
		}
		teamID = &parsed
	}

	user, err := h.service.SetTeam(r.Context(), targetID, teamID)
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toUserResponse(user))
}
