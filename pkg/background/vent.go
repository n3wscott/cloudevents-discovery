package background

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/types"
	"log"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/subscription"
)

type ServiceChange struct {
	Change  string            `json:"change"`
	Service discovery.Service `json:"service"`
}

type SubscriptionChange struct {
	Change       string                    `json:"change"`
	Subscription subscription.Subscription `json:"subscription"`
}

func NewVent(service string, sinks string, changes <-chan ServiceChange, subs <-chan SubscriptionChange) Background {
	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	sl := make([]types.URI, 0)
	for _, s := range strings.Split(sinks, ",") {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		if u := types.ParseURI(s); u != nil {
			sl = append(sl, *u)
		}
	}

	return &Vent{
		service: service,
		changes: changes,
		subs:    subs,
		client:  client,
		sinks:   sl,
	}
}

type Vent struct {
	service string
	changes <-chan ServiceChange
	subs    <-chan SubscriptionChange

	client cloudevents.Client
	sinks  []types.URI
}

func (v *Vent) eventFor(change ServiceChange) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetType(fmt.Sprintf("cloudmeta.discovery.service.%s.v1", change.Change))
	event.SetSource(v.service)
	event.SetSubject(fmt.Sprintf("/services/%s", change.Service.ID))
	if err := event.SetData(cloudevents.ApplicationJSON, &change); err != nil {
		return nil, err
	}
	return &event, nil
}

func (v *Vent) Start(ctx context.Context) error {
	for {
		select {
		case change := <-v.subs:
			fmt.Println("---------------------")
			fmt.Println(change.Change, change.Subscription.Sink)
			switch change.Change {
			case "added":
				// TODO: we do not use the filters yet...
				v.sinks = append(v.sinks, change.Subscription.Sink)

				ctx := cloudevents.ContextWithTarget(context.Background(), change.Subscription.Sink.String())

				event := cloudevents.NewEvent()
				event.SetType("cloudmeta.discovery.service.subscribed.v1")
				event.SetSource(v.service)
				event.SetSubject("/subscriptions/" + change.Subscription.ID)

				if result := v.client.Send(ctx, event); cloudevents.IsUndelivered(result) {
					log.Println("failed to deliver event to sink: ", change.Subscription.Sink)
				}

			case "updated":
				log.Println("todo: updated")
			case "deleted":
				log.Println("todo: deleted")
			}
			fmt.Println("---------------------")

		case change := <-v.changes:
			fmt.Println("service updated:", change.Change, "for", change.Service.Name)

			for _, sink := range v.sinks {
				ctx := cloudevents.ContextWithTarget(context.Background(), sink.String())

				event, err := v.eventFor(change)
				if err != nil {
					log.Println("failed to create event for change: ", change)
					continue
				}

				if result := v.client.Send(ctx, *event); cloudevents.IsUndelivered(result) {
					log.Println("failed to deliver event to sink: ", sink)
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
