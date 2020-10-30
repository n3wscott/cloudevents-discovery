package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/cloudevents-discovery/pkg/background"
	"github.com/n3wscott/cloudevents-discovery/pkg/handler"
	"log"
	"net/http"
	"os"
)

type envConfig struct {
	Service       string `envconfig:"SERVICE" default:"http://localhost:8080"`
	Port          int    `envconfig:"PORT" default:"8080"`
	Downstream    string `envconfig:"DISCOVERY_DOWNSTREAM"` // comma separated list of urls.
	Services      string `envconfig:"DISCOVERY_SERVICES_FILE"`
	Subscriptions string `envconfig:"SUBSCRIPTIONS_FILE"`
	Sinks         string `envconfig:"SINK"` // comma separated list of urls.
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	changes := make(chan background.ServiceChange, 10)   // TODO: 10 might be too small of a channel buffer.
	subs := make(chan background.SubscriptionChange, 10) // TODO: 10 might be too small of a channel buffer.

	servicesHandler := handler.NewServiceHandler(changes)
	if env.Services != "" {
		if err := servicesHandler.LoadServicesFromFile(env.Services); err != nil {
			log.Fatal(err)
		}
	}
	// Add ourself.
	servicesHandler.Set(background.Service(env.Service))

	subscriptionHandler := handler.NewSubscriptionHandler(subs)

	r := mux.NewRouter()

	r.Handle("/services", servicesHandler)
	r.Handle("/services/{id}", servicesHandler)

	r.Handle("/subscriptions", subscriptionHandler)
	r.Handle("/subscriptions/{id}", subscriptionHandler)

	http.Handle("/", r)

	ctx := context.Background()

	vent := background.NewVent(env.Service, env.Sinks, changes, subs)
	go func() {
		if err := vent.Start(ctx); err != nil {
			log.Println(err)
		}
	}()

	agg := background.NewDiscoveryAggregation(env.Downstream, servicesHandler)
	go func() {
		if err := agg.Start(ctx); err != nil {
			log.Println(err)
		}
	}()

	addr := fmt.Sprintf(":%d", env.Port)

	log.Printf("will listen on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
