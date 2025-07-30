package server

import (
	"encoding/json"
	"github.com/AndreySirin/04.08/internal/entity"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const (
	statusWaiting = "waiting"
	statusload    = "loading"
)

func (s *Server) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task string
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var slise []string
	s.Task[id] = &entity.Task{
		Link:   slise,
		Status: statusWaiting,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s.Task)
}

func (s *Server) HandleCreateLink(w http.ResponseWriter, r *http.Request) {
	var url string
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(s.Task[id].Link) == 2 {
		s.Task[id].Link = append(s.Task[id].Link, url)
		s.Task[id].Status = statusload

	} else if len(s.Task[id].Link) < 3 {
		s.Task[id].Link = append(s.Task[id].Link, url)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s.Task[id])
}

func (s *Server) HandleGetStatus(w http.ResponseWriter, r *http.Request) {

	taskId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.Task[id].Status != statusload {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(s.Task[id].Status)
	} else {

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("надо загружать и архивировать")
	}
}
