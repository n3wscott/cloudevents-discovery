package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
)

type ServicesHandler struct {
	once     sync.Once
	services []discovery.Service
}

func (h *ServicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		h.services = make([]discovery.Service, 0)
		err := json.Unmarshal([]byte(exampleServices), &h.services)
		if err != nil {
			panic(err)
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

func (h *ServicesHandler) handleList(w http.ResponseWriter, r *http.Request) {
	services := h.services

	name := r.URL.Query().Get("name")
	if name != "" {
		name = strings.ToLower(name)
		services = make([]discovery.Service, 0)
		for _, v := range h.services {
			if strings.ToLower(v.Name) == name {
				services = append(services, v)
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

func (h *ServicesHandler) handleGet(id string, w http.ResponseWriter, r *http.Request) {
	var service *discovery.Service

	for _, v := range h.services {
		if v.ID == id {
			service = &v
		}
	}

	if service == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	js, err := json.Marshal(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
