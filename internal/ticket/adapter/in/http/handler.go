package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
	"github.com/jhonsferg/corelog/internal/platform/httpresponse"
	appmw "github.com/jhonsferg/corelog/internal/platform/middleware"
	"github.com/jhonsferg/corelog/internal/platform/validation"
	"github.com/jhonsferg/corelog/internal/ticket/application/port"
	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type Handler struct {
	service   port.Service
	jwtSecret string
}

func NewHandler(service port.Service, jwtSecret string) *Handler {
	return &Handler{service: service, jwtSecret: jwtSecret}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := appmw.UserIDFromContext(r.Context())
	if !ok {
		httpresponse.Error(w, apperror.New(apperror.KindUnauthorized, "not authenticated"))
		return
	}

	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	ticket, err := h.service.Create(r.Context(), port.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    domain.Priority(req.Priority),
		RequesterID: requesterID,
	})
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusCreated, toTicketResponse(ticket))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid ticket id"))
		return
	}

	ticket, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toTicketResponse(ticket))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := port.ListFilter{
		Status:   domain.Status(q.Get("status")),
		Priority: domain.Priority(q.Get("priority")),
		Page:     atoiOrDefault(q.Get("page"), 1),
		PageSize: atoiOrDefault(q.Get("page_size"), 20),
	}

	tickets, err := h.service.List(r.Context(), filter)
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toTicketResponses(tickets))
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid ticket id"))
		return
	}

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	ticket, err := h.service.UpdateStatus(r.Context(), id, domain.Status(req.Status))
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toTicketResponse(ticket))
}

func (h *Handler) Assign(w http.ResponseWriter, r *http.Request) {
	assignerID, ok := appmw.UserIDFromContext(r.Context())
	if !ok {
		httpresponse.Error(w, apperror.New(apperror.KindUnauthorized, "not authenticated"))
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid ticket id"))
		return
	}

	var req assignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid request body"))
		return
	}
	if err := validation.Struct(req); err != nil {
		httpresponse.Error(w, apperror.Wrap(apperror.KindValidation, "invalid request body", err))
		return
	}

	assigneeID, err := uuid.Parse(req.AssigneeID)
	if err != nil {
		httpresponse.Error(w, apperror.New(apperror.KindValidation, "invalid assignee id"))
		return
	}

	ticket, err := h.service.Assign(r.Context(), id, assignerID, assigneeID)
	if err != nil {
		httpresponse.Error(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, toTicketResponse(ticket))
}

func atoiOrDefault(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}
