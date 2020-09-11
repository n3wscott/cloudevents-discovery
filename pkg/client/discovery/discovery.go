package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
)

type DiscoveryAPI interface {
	Services() Services
}

type Services interface {
	Get(ctx context.Context, id string, opts *GetOptions) (*discovery.Service, error)
	List(ctx context.Context, opts *ListOptions) ([]discovery.Service, error)
}

type GetOptions struct {
}

type ListOptions struct {
	Name string
}

// client.Discovery("url").Services().Get(id)
// client.Discovery("url").Services().List(opts)

func New(baseURL url.URL) DiscoveryAPI {
	return &client{baseURL: baseURL}
}

type client struct {
	baseURL url.URL
}

func (c *client) Services() Services {
	return &services{c: c}
}

type services struct {
	c *client
}

func (s *services) Get(ctx context.Context, id string, _ *GetOptions) (*discovery.Service, error) {
	target := fmt.Sprintf("%s/services/%s", s.c.baseURL.String(), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d, %s", resp.Body, string(b))
	}

	svc := new(discovery.Service)
	if err := json.NewDecoder(resp.Body).Decode(svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *services) List(ctx context.Context, opts *ListOptions) ([]discovery.Service, error) {
	target := fmt.Sprintf("%s/services", s.c.baseURL.String())
	if opts != nil && opts.Name != "" {
		target = fmt.Sprintf("%s?name=%s", target, opts.Name)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d, %s", resp.Body, string(b))
	}

	svcs := make([]discovery.Service, 0)
	if err := json.NewDecoder(resp.Body).Decode(&svcs); err != nil {
		return nil, err
	}
	return svcs, nil
}
