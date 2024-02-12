package internal

/*
Copyright © 2024 Pete Wall <pete@petewall.net>
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

const DefaultPort = 8080

type Server struct {
	Port      int
	Datastore Datastore
}

func (s *Server) Start() error {
	log.Infof("Starting HTTP server on port %d...", s.Port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	r.Put("/api/okr", s.addOKR)
	r.Post("/api/okr", s.updateOKR)
	r.Get("/api/okr/{id}", s.getOKR)
	r.Delete("/api/okr/{id}", s.deleteOKR)

	r.Get("/api/okrs", s.getAllOKRs)
	r.Get("/api/okrs/{quarter}", s.getOKRsByQuarter)

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), r)
}

func (s *Server) getOKR(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	okr, err := s.Datastore.Get(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("failed to get OKR")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to get OKR %s", id)
		return
	}

	if okr == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "OKR with id %s not found", id)
		return
	}

	data, err := json.Marshal(okr)
	if err != nil {
		log.WithError(err).Error("failed to convert OKR into JSON")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to convert OKR into JSON")
		return
	}
	_, _ = w.Write(data)
}

func (s *Server) addOKR(w http.ResponseWriter, r *http.Request) {
	var okr *OKR
	err := json.NewDecoder(r.Body).Decode(&okr)
	if err != nil {
		log.WithError(err).Error("failed to parse OKR")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "failed to parse OKR: %s", err.Error())
		return
	}

	s.Datastore.Add(CreateOKR(okr.Quarter, okr.Category, okr.ValueType, okr.Description, okr.Goal))
}

func (s *Server) updateOKR(w http.ResponseWriter, r *http.Request) {
	var okr *OKR
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "failed to read request body: %s", err.Error())
		return
	}

	err = json.Unmarshal(body, &okr)
	if err != nil {
		log.WithError(err).WithField("body", string(body)).Error("failed to parse OKR")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "failed to parse OKR: %s", err.Error())
		return
	}

	s.Datastore.Update(okr)
}

func (s *Server) deleteOKR(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := s.Datastore.Delete(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("failed to delete OKR")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to delete OKR %s", id)
		return
	}
	_, _ = fmt.Fprintf(w, "OK")
}

func (s *Server) getAllOKRs(w http.ResponseWriter, r *http.Request) {
	okrs, _ := s.Datastore.GetAll()
	data, err := json.Marshal(okrs)
	if err != nil {
		log.WithError(err).Error("failed to convert OKRs into JSON")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to convert OKRs into JSON")
		return
	}
	_, _ = w.Write(data)
}

func (s *Server) getOKRsByQuarter(w http.ResponseWriter, r *http.Request) {
	quarter := chi.URLParam(r, "quarter")
	okrs, _ := s.Datastore.GetByQuarter(quarter)
	data, err := json.Marshal(okrs)
	if err != nil {
		log.WithError(err).Error("failed to convert OKRs into JSON")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to convert OKRs into JSON")
		return
	}
	_, _ = w.Write(data)
}
