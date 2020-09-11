package background

import (
	"context"
	"net/url"
	"strings"
)

type DiscoveryAggregation struct {
	DownStream []url.URL
}

// downstream is a comma separated list of urls.
func NewDiscoveryAggregation(downstream string) Background {
	ds := make([]url.URL, 0)
	for _, s := range strings.Split(downstream, ",") {
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		ds = append(ds, *u)
	}
	return &DiscoveryAggregation{DownStream: ds}
}

func (a *DiscoveryAggregation) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
