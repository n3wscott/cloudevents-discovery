package subscription

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/n3wscott/cloudevents-discovery/pkg/apis/subscription"
)

type SubscriptionAPI interface {
	Subscriptions() Subscription
}

type Subscription interface {
	Create(ctx context.Context, s subscription.Subscription, opts *CreateOptions) (*subscription.Subscription, error)
	Update(ctx context.Context, s subscription.Subscription, opts *UpdateOptions) (*subscription.Subscription, error)
	Delete(ctx context.Context, id string, opts *DeleteOptions) error
	Get(ctx context.Context, id string, opts *GetOptions) (*subscription.Subscription, error)
	List(ctx context.Context, opts *ListOptions) ([]subscription.Subscription, error)
}

type CreateOptions struct{}
type UpdateOptions struct{}
type DeleteOptions struct{}
type GetOptions struct{}

type ListOptions struct {
	Name string
}

// client.Subscription("url").Subscriptions().Get(id)
// client.Subscription("url").Subscriptions().List(opts)

func New(baseURL url.URL) SubscriptionAPI {
	return &client{baseURL: baseURL}
}

type client struct {
	baseURL url.URL
}

func (c *client) Subscriptions() Subscription {
	return &subscriptions{c: c}
}

type subscriptions struct {
	c *client
}

func (s *subscriptions) Create(ctx context.Context, up subscription.Subscription, _ *CreateOptions) (*subscription.Subscription, error) {
	target := fmt.Sprintf("%s/subscriptions", s.c.baseURL.String())

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(up); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, b)
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

	sub := new(subscription.Subscription)
	if err := json.NewDecoder(resp.Body).Decode(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *subscriptions) Update(ctx context.Context, up subscription.Subscription, _ *UpdateOptions) (*subscription.Subscription, error) {
	target := fmt.Sprintf("%s/subscriptions", s.c.baseURL.String())

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(up); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, target, b)
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

	sub := new(subscription.Subscription)
	if err := json.NewDecoder(resp.Body).Decode(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *subscriptions) Delete(ctx context.Context, id string, _ *DeleteOptions) error {
	target := fmt.Sprintf("%s/subscriptions/%s", s.c.baseURL.String(), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, target, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%d, %s", resp.Body, string(b))
	}
	return nil
}

func (s *subscriptions) Get(ctx context.Context, id string, _ *GetOptions) (*subscription.Subscription, error) {
	target := fmt.Sprintf("%s/subscriptions/%s", s.c.baseURL.String(), id)
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

	sub := new(subscription.Subscription)
	if err := json.NewDecoder(resp.Body).Decode(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *subscriptions) List(ctx context.Context, _ *ListOptions) ([]subscription.Subscription, error) {
	target := fmt.Sprintf("%s/subscriptions", s.c.baseURL.String())

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

	subs := make([]subscription.Subscription, 0)
	if err := json.NewDecoder(resp.Body).Decode(&subs); err != nil {
		return nil, err
	}
	return subs, nil
}
