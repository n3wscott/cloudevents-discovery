package background

import (
	"context"
	"fmt"
	"github.com/n3wscott/cloudevents-discovery/pkg/client"
	"net/url"
	"strings"
	"time"
)

type discoveryAggregation struct {
	downstream []url.URL
	period     time.Duration
	mgr        ServicesManager
}

// downstream is a comma separated list of urls.
func NewDiscoveryAggregation(downstream string, mgr ServicesManager) Background {
	ds := make([]url.URL, 0)
	for _, s := range strings.Split(downstream, ",") {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		ds = append(ds, *u)
	}
	return &discoveryAggregation{
		downstream: ds,
		period:     time.Second * 10,
		mgr:        mgr,
	}
}

func (a *discoveryAggregation) Start(ctx context.Context) error {
	c := client.New()
	timer := time.Tick(a.period)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Discovery Aggregation - done")
			return nil
		case <-timer:
			fmt.Printf("Discovery.Services[%d]\n", len(a.downstream))

			for _, d := range a.downstream {
				fmt.Printf("\tfor %s,\n", d.String())
				svcs, err := c.Discovery(d).Services().List(ctx, nil)
				if err != nil {
					fmt.Printf("\tfailed to list services from %s, %v", d.String(), err)
					continue
				}
				for i, svc := range svcs {
					fmt.Printf("\t[%d]:\t%+v\n", i, svc)
					a.mgr.Set(svc)
				}
			}
		}
	}
}
