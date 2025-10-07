package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"movieexample.com/rating/internal/controller"
	model "movieexample.com/rating/pkg"
)

// Handler define a rating service controller
type Handler struct {
	ctrl *controller.Controller
}

// New creates a new rating service HTTP handler.
func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// PutRating create a service to put ratings in records
func (h *Handler) PutRating(w http.ResponseWriter, req *http.Request) {
	recordID := model.RecordID(req.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(req.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := model.UserID(req.FormValue("userId"))
	v, err := strconv.ParseFloat(req.FormValue("value"),
		64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rating := &model.Rating{
		UserID: userID,
		Value:  model.RatingValue(v),
	}

	if err := h.ctrl.PutRating(req.Context(), recordID, recordType, rating); err != nil {
		log.Printf("Repository put error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetRatings retrives the rating aggregated
func (h *Handler) GetRatings(w http.ResponseWriter, req *http.Request) {

	recordID := model.RecordID(req.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(req.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, err := h.ctrl.GetAggregatedRating(req.Context(), recordID, recordType)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}

}
