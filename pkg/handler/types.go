package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
)

type TypesHandler struct {
	once     sync.Once
	services map[string][]discovery.Service
}

func (h *TypesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		h.services = make(map[string][]discovery.Service, 0)
		services := make([]discovery.Service, 0)
		err := json.Unmarshal([]byte(exampleServices), &services)
		if err != nil {
			panic(err)
		}
		for _, service := range services {
			for _, t := range service.Types {
				if _, found := h.services[t.Type]; !found {
					h.services[t.Type] = []discovery.Service{service}
				} else {
					h.services[t.Type] = append(h.services[t.Type], service)
				}
			}
		}
	})

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		h.handleList(w, r)
	} else {
		h.handleGet(id, w, r)
	}
}

func (h *TypesHandler) handleList(w http.ResponseWriter, r *http.Request) {
	services := h.services

	matching := r.URL.Query().Get("matching")
	if matching != "" {
		matching = strings.ToLower(matching)
		services = make(map[string][]discovery.Service, 0)
		for k, v := range h.services {
			if strings.Contains(strings.ToLower(k), matching) {
				services[k] = v
			}
		}
	}

	js, err := json.Marshal(services)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (h *TypesHandler) handleGet(id string, w http.ResponseWriter, r *http.Request) {
	if _, found := h.services[id]; !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	services := map[string][]discovery.Service{
		id: h.services[id],
	}

	js, err := json.Marshal(services)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
