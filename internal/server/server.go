package server

import (
	"context"
	"errors"
	"github.com/AndreySirin/04.08/internal/client"
	"github.com/AndreySirin/04.08/internal/entity"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	HttpServer *http.Server
	Task       map[int]*entity.Task
	ch         chan string
	lg         *slog.Logger
	client     *http.Client
}

func New(addr string, logger *slog.Logger) *Server {
	m := make(map[int]*entity.Task)
	ch := make(chan string, 3)
	cl := client.New()
	s := &Server{
		Task:   m,
		lg:     logger,
		client: cl,
		ch:     ch,
	}

	r := chi.NewRouter()
	r.Get("/archives/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/archives/", http.FileServer(http.Dir("archives")))
		fs.ServeHTTP(w, r)
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/tasks", s.HandleCreateTask)
			r.Post("/links/{id}", s.HandleCreateLink)
			r.Get("/status/{id}", s.HandleGetStatus)
		})
	})
	s.HttpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return s
}
func (s *Server) Run() {
	s.lg.Info("server is running")
	err := s.HttpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		s.lg.Info("server is stopped")
		return
	}
}

func (s *Server) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.HttpServer.Shutdown(ctx)
	if err != nil {
		s.lg.Error("error when stopping the server")
		return
	}
}
