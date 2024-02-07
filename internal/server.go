package internal

/*
Copyright © 2024 Pete Wall <pete@petewall.net>
*/

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Port      int
	Datastore Datastore
}

func (s *Server) Start() error {
	log.Info("Starting HTTP server...")
	log.Debugf("Using port %d", s.Port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/api/okrs", s.handleGetAllOKRs)
	r.Get("/api/okrs/<quarter", promhttp.Handler().ServeHTTP)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), r)
}

func (s *Server) handleGetAllOKRs(w http.ResponseWriter, r *http.Request) {
	okrs := s.Datastore.GetAll()
	data, _ := json.Marshal(okrs)
	_, _ = w.Write(data)
}

func (s *Server) handleGetOKRsByQuarter(w http.ResponseWriter, r *http.Request) {
	quarter := chi.URLParam(r, "quarter")
	okrs := s.Datastore.GetByQuarter(quarter)
	data, _ := json.Marshal(okrs)
	_, _ = w.Write(data)
}
