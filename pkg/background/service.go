package background

import (
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/discovery"
	"net/url"
	"strings"
)

func Service(service string) discovery.Service {
	host := "failed"
	if u, err := url.Parse(service); err == nil {
		host = strings.ReplaceAll(u.Host, ":", "")
	}

	id := host + "-todo-123"
	return discovery.Service{
		ID:              id,
		URL:             service + "/services/" + id,
		Name:            "CloudMeta",
		Epoch:           0,
		Description:     "Discovery Service Events from CloudMeta",
		SpecVersions:    []string{"1.0"},
		SubscriptionURL: service + "/subscriptions",
		Protocols:       []string{"HTTP"},
		Events: []discovery.ServiceEvent{{
			Type:            "cloudmeta.discovery.service.subscribed.v1",
			Description:     "Discovery - Service entry subscription start of stream.",
			DataContentType: "application/json",
		}, {
			Type:            "cloudmeta.discovery.service.added.v1",
			Description:     "Discovery - Service entry was added.",
			DataContentType: "application/json",
		}, {
			Type:            "cloudmeta.discovery.service.updated.v1",
			Description:     "Discovery - Service entry was updated.",
			DataContentType: "application/json",
		}, {
			Type:            "cloudmeta.discovery.service.deleted.v1",
			Description:     "Discovery - Service entry was deleted.",
			DataContentType: "application/json",
		}},
	}
}
