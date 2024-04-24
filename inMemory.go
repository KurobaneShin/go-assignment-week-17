package main

import (
	"encoding/json"
	"net/http"
)

type MoResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *Server) inMemoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleMemoryObjectGet(w, r)
	case http.MethodPost:
		s.handleMemoryObjectPost(w, r)
	default:
		http.Error(w, "Method Not Implemented", http.StatusNotImplemented)
	}
}

func (s *Server) handleMemoryObjectGet(w http.ResponseWriter, r *http.Request) {
	keyParam := r.URL.Query().Get("key")
	if keyParam == "" {
		http.Error(w, "missing query param: key", http.StatusBadRequest)
		return
	}
	val, ok := s.mo.Get(keyParam)

	if !ok {
		http.Error(w, "value not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&MoResponse{
		Key:   keyParam,
		Value: val,
	})
}

func (s *Server) handleMemoryObjectPost(w http.ResponseWriter, r *http.Request) {
	var object MoResponse
	err := json.NewDecoder(r.Body).Decode(&object)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	s.mo.Set(object.Key, object.Value)
	w.WriteHeader(http.StatusCreated)
}
