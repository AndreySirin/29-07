package server

import (
	"encoding/json"
	"github.com/AndreySirin/04.08/internal/archive"
	"github.com/AndreySirin/04.08/internal/client"
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
		writer, file, errorZip := archive.CreateZip(taskId)
		defer file.Close()
		defer writer.Close()
		if errorZip != nil {
			http.Error(w, errorZip.Error(), http.StatusBadRequest)
			return
		}
		for i := 0; i < 3; i++ {
			body, errLoad := client.LoadFile(s.Task[id].Link[i], s.client)
			if errLoad != nil {
				http.Error(w, errLoad.Error(), http.StatusBadRequest)
				return
			}
			// проверка content-type
			err = archive.WriteZip(writer, body, s.Task[id].Link[i])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode("надо загружать и архивировать")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
