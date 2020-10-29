package background

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
)

type ServiceChange struct {
	Change  string            `json:"change"`
	Service discovery.Service `json:"service"`
}

func NewVent(sinks string, changes <-chan ServiceChange) Background {

	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	sl := make([]url.URL, 0)
	for _, s := range strings.Split(sinks, ",") {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		u, err := url.Parse(s)
		if err != nil {
			log.Printf("invalid sink, skipping %s\n", s)
			continue
		}
		sl = append(sl, *u)
	}

	return &Vent{
		changes: changes,
		client:  client,
		sinks:   sl,
	}
}

type Vent struct {
	changes <-chan ServiceChange

	client cloudevents.Client
	sinks  []url.URL
}

func eventFor(change ServiceChange) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetType(fmt.Sprintf("n3wscott.discovery.service.%s.v1", change.Change))
	event.SetSource("todo://host")
	event.SetSubject(fmt.Sprintf("/services/%s", change.Service.ID))
	if err := event.SetData(cloudevents.ApplicationJSON, &change); err != nil {
		return nil, err
	}
	return &event, nil
}

func (v *Vent) Start(ctx context.Context) error {
	for {
		select {
		case change := <-v.changes:
			fmt.Println(change)

			for _, sink := range v.sinks {
				ctx := cloudevents.ContextWithTarget(context.Background(), sink.String())

				event, err := eventFor(change)
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
