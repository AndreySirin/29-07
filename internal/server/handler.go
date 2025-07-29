package server

import (
	"encoding/json"
	"github.com/AndreySirin/04.08/internal/entity"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	statusWaiting = "waiting"
)

func (s *Server) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var str string
	if err := json.NewDecoder(r.Body).Decode(&str); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.Task[str] = &entity.Task{
		Status: statusWaiting,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s.Task)
}

func (s *Server) HandleCreateLink(w http.ResponseWriter, r *http.Request) {

	var str string
	if err := json.NewDecoder(r.Body).Decode(&str); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := chi.URLParam(r, "task")
	str = s.Task[task].Link[0]
	s.Task[task] = &entity.Task{
		Link:   s.Task[task].Link,
		Status: statusWaiting,
	}

}
