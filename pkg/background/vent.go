package background

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/subscription"
	"log"
	"strings"
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

	sl := make([]subscription.Subscription, 0)
	for i, s := range strings.Split(sinks, ",") {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		if u := types.ParseURI(s); u != nil {
			sl = append(sl, subscription.Subscription{
				ID:       fmt.Sprintf("manual-entry-%d", i),
				Protocol: "HTTP",
				Sink:     *u,
			})
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
	sinks  []subscription.Subscription
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
				v.sinks = append(v.sinks, change.Subscription)

				ctx := cloudevents.ContextWithTarget(context.Background(), change.Subscription.Sink.String())

				event := cloudevents.NewEvent()
				event.SetType("cloudmeta.discovery.service.subscribed.v1")
				event.SetSource(v.service)
				event.SetSubject("/subscriptions/" + change.Subscription.ID)

				if change.Subscription.Filter != nil && basicFiltered(&event, change.Subscription.Filter.Filters) {
					continue
				}
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

			event, err := v.eventFor(change)
			if err != nil {
				log.Println("failed to create event for change: ", change)
				continue
			}

			for _, sink := range v.sinks {

				if sink.Protocol != "HTTP" {
					log.Println("skipping: Subscription not HTTP", sink.ID)
					continue
				}

				if sink.Filter != nil {
					if sink.Filter.Dialect != "basic" {
						log.Println("skipping: Subscription filter not supported, ", sink.ID, sink.Filter.Dialect)
						continue
					}

					if basicFiltered(event, sink.Filter.Filters) {
						continue
					}
				}

				ctx := cloudevents.ContextWithTarget(context.Background(), sink.Sink.String())
				if result := v.client.Send(ctx, *event); cloudevents.IsUndelivered(result) {
					log.Println("failed to deliver event to sink: ", sink)
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func basicFiltered(event *cloudevents.Event, filters []subscription.BasicFilter) bool {
	for _, f := range filters {
		value := ""
		switch f.Property {
		case "specversion":
			value = event.SpecVersion()
		case "type":
			value = event.Type()
		case "source":
			value = event.Source()
		case "subject":
			value = event.Subject()
		case "id":
			value = event.ID()
		case "time":
			value, _ = types.ToString(event.Time())
		case "dataschema":
			value = event.DataSchema()
		case "datacontenttype":
			value = event.DataContentType()
		default:
			value, _ = types.ToString(event.Extensions()[f.Property])
		}

		switch f.Type {
		case "prefix":
			if !strings.HasPrefix(value, f.Value) {
				return true
			}
		case "suffix":
			if !strings.HasSuffix(value, f.Value) {
				return true
			}
		case "exact":
			if value != f.Value {
				return true
			}
		default:
			return true
		}
	}
	return false
}
