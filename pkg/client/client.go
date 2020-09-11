package client

import (
	"net/url"

	"github.com/n3wscott/cloudevents-discovery/pkg/client/discovery"
	"github.com/n3wscott/cloudevents-discovery/pkg/client/subscription"
)

type Client interface {
	Subscriptions(baseURL url.URL) subscription.SubscriptionAPI
	Discovery(baseURL url.URL) discovery.DiscoveryAPI
}

func New() Client {
	return &client{}
}

type client struct {
	// TODO: this will hold http clients and auth and stuff.
}

func (c *client) Subscriptions(baseURL url.URL) subscription.SubscriptionAPI {
	return subscription.New(baseURL)
}

func (c *client) Discovery(baseURL url.URL) discovery.DiscoveryAPI {
	return discovery.New(baseURL)
}
