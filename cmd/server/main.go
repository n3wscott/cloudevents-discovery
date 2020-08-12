package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/n3wscott/cloudevents-discovery/pkg/handler"
)

func main() {

	servicesHandler := new(handler.ServicesHandler)
	typesHandler := new(handler.TypesHandler)

	r := mux.NewRouter()

	r.Handle("/services", servicesHandler)
	r.Handle("/services/{id}", servicesHandler)
	r.Handle("/types", typesHandler)
	r.Handle("/types/{id}", typesHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
