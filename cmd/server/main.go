package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/n3wscott/cloudevents-discovery/pkg/handler"
)

func main() {
	servicesHandler := new(handler.ServicesHandler)
	subscriptionHandler := new(handler.SubscriptionHandler)

	r := mux.NewRouter()

	r.Handle("/services", servicesHandler)
	r.Handle("/services/{id}", servicesHandler)

	r.Handle("/subscriptions", subscriptionHandler)
	r.Handle("/subscriptions/{id}", subscriptionHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
