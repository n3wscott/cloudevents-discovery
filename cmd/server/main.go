package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/cloudevents-discovery/pkg/background"
	"log"
	"net/http"
	"os"

	"github.com/n3wscott/cloudevents-discovery/pkg/handler"
)

type envConfig struct {
	Port          int    `envconfig:"PORT" default:"8080"`
	Downstream    string `envconfig:"DISCOVERY_DOWNSTREAM"` // comma separated list of urls.
	Services      string `envconfig:"DISCOVERY_SERVICES_FILE"`
	Subscriptions string `envconfig:"SUBSCRIPTIONS_FILE"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	servicesHandler := new(handler.ServicesHandler)
	if env.Services != "" {
		if err := servicesHandler.LoadServicesFromFile(env.Services); err != nil {
			log.Fatal(err)
		}
	}

	subscriptionHandler := new(handler.SubscriptionHandler)

	r := mux.NewRouter()

	r.Handle("/services", servicesHandler)
	r.Handle("/services/{id}", servicesHandler)

	r.Handle("/subscriptions", subscriptionHandler)
	r.Handle("/subscriptions/{id}", subscriptionHandler)

	http.Handle("/", r)

	ctx := context.Background()

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
