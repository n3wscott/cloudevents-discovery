package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/subscription"
	"github.com/n3wscott/cloudevents-discovery/pkg/client"
)

type envConfig struct {
	Discovery    string `envconfig:"DISCOVERY_HOST" default:"http://localhost:8080"`
	Subscription string `envconfig:"SUBSCRIPTION_HOST" default:"http://localhost:8080"`
}

func main() {
	ctx := context.Background()
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	d, err := url.Parse(env.Discovery)
	if err != nil {
		panic(err)
	}

	s, err := url.Parse(env.Subscription)
	if err != nil {
		panic(err)
	}

	c := client.New()

	fmt.Printf("Discovery.Services:\n")

	svcs, err := c.Discovery(*d).Services().List(ctx, nil)
	if err != nil {
		panic(err)
	}
	for i, svc := range svcs {
		fmt.Printf("[%d]:\t%+v\n", i, svc)
	}

	fmt.Printf("Subscriptions.Subscriptions:\n")

	subs, err := c.Subscriptions(*s).Subscriptions().List(ctx, nil)
	if err != nil {
		panic(err)
	}
	for i, sub := range subs {
		fmt.Printf("[%d]:\t%+v\n", i, sub)
	}

	sink := cloudevents.ParseURI("https://localhost:8080")
	sub := subscription.Subscription{
		ID:       "123-123",
		Protocol: "HTTP",
		Sink:     *sink,
	}

	// Uncomment for filtered example.
	//sub.Filter = &subscription.Filter{
	//	Dialect: "basic",
	//	Filters: []subscription.BasicFilter{{
	//		Type:     "suffix",
	//		Property: "type",
	//		Value:    "updated.v1",
	//	}},
	//}

	updated, err := c.Subscriptions(*s).Subscriptions().Update(ctx, sub, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[updated]:\t%+v\n", updated)
}
