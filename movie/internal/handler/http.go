package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieexample.com/movie/internal/controller"
)

// Handler defines a movie handler.
type Handler struct {
	ctrl *controller.Controller
}

// New creates a new movie HTTP handler.
func New(ctrl *controller.Controller) *Handler {

	return &Handler{
		ctrl: ctrl,
	}
}

// GetMovieDetails hanldes GET/movie request
func (h *Handler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return

	}

	ctx := r.Context()
	movieDetails, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal server error :%v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(movieDetails); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
