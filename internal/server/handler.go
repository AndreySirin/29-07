package server

import (
	"encoding/json"
	"fmt"
	"github.com/AndreySirin/04.08/internal/client"
	"github.com/AndreySirin/04.08/internal/entity"
	"github.com/AndreySirin/04.08/internal/zip"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
)

const (
	statusCreated   = "created"
	statusLoad      = "loading"
	statusCompleted = "completed"
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
	select {
	case s.ch <- task:
		var links []string
		er := make(map[string]string)
		s.Task[id] = &entity.Task{
			Link:      links,
			Status:    statusCreated,
			ErrorLoad: er,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(s.Task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("the server is currently busy")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
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
		s.Task[id].Status = statusLoad

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
	if s.Task[id].Status == statusCreated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(s.Task[id].Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if s.Task[id].Status == statusLoad {
		writer, file, errorZip := zip.CreateZip(taskId)
		if errorZip != nil {
			http.Error(w, errorZip.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			err = file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		defer func() {
			err = writer.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		for i := 0; i < 3; i++ {
			body, errLoad := client.LoadFile(s.Task[id].Link[i], s.client)
			if errLoad != nil {
				s.Task[id].ErrorLoad[s.Task[id].Link[i]] = errLoad.Error()
				s.lg.Error("error when uploading a file from a link", slog.String("url", s.Task[id].Link[i]))
				continue
			}
			err = zip.WriteZip(writer, body, s.Task[id].Link[i])
			if err != nil {
				s.lg.Error("error when writing to the zip", err)
				continue
			}
		}
		linkZip := fmt.Sprintf("http://localhost:8080/archives/%s.zip", taskId)
		s.Task[id].ArchiveUrl = linkZip
		s.Task[id].Status = statusCompleted
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(s.Task[id])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		<-s.ch
	} else {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("the server has already completed this task")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
