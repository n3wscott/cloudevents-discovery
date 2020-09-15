package background

import (
	"context"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
)

type Background interface {
	Start(context.Context) error
}

type ServicesManager interface {
	Set(service discovery.Service)
}
