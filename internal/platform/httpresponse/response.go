package httpresponse

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
)

type envelope struct {
	Data any `json:"data"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(envelope{Data: data}); err != nil {
		slog.Error("httpresponse: encode failed", "error", err)
	}
}

func Error(w http.ResponseWriter, err error) {
	status := statusFor(apperror.KindOf(err))
	if status == http.StatusInternalServerError {
		slog.Error("unhandled error", "error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: errorBody{Message: safeMessage(status, err)}})
}

func statusFor(kind apperror.Kind) int {
	switch kind {
	case apperror.KindValidation:
		return http.StatusUnprocessableEntity
	case apperror.KindNotFound:
		return http.StatusNotFound
	case apperror.KindConflict:
		return http.StatusConflict
	case apperror.KindUnauthorized:
		return http.StatusUnauthorized
	case apperror.KindForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func safeMessage(status int, err error) string {
	if status == http.StatusInternalServerError {
		return "an unexpected error occurred"
	}
	return err.Error()
}
