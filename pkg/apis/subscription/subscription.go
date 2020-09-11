package subscription

import (
	"bytes"
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Based on https://github.com/cloudevents/spec/tree/e9f516548f369f15348d4f1603bbc98c342a1d8f/subscriptions-api.md

type Subscription struct {
	// ID - The unique identifier of the subscription in the scope of the subscription manager.
	ID string `json:"id"`

	// Protocol - Identifier of a delivery protocol. Because of WebSocket tunneling options for AMQP, MQTT and other
	// protocols, the URI scheme is not sufficient to identify the protocol. The protocols with existing CloudEvents
	// bindings are identified as "AMQP", "MQTT3", "MQTT5", "HTTP", "KAFKA", and "NATS". An implementation MAY add
	// support for further protocols.
	Protocol string `json:"protocol"`

	// ProtocolSettings - A set of settings specific to the selected delivery protocol provider. Options for those
	// settings are listed in the following subsection. An implementation MAY offer more options. Examples for such
	// settings are credentials, which generally vary by transport, rate limits and retry policies, or the QoS mode for
	// MQTT. See the Protocol Settings section for further details.
	// +optional
	ProtocolSettings *json.RawMessage `json:"protocolsettings,omitempty"`

	// Sink - The address to which events SHALL be delivered using the selected protocol. The format of the address MUST
	// be valid for the selected protocol or one of the protocol's own transport bindings (e.g. AMQP over WebSockets).
	Sink cloudevents.URIRef `json:"sink"`

	/// Filter - A filter is an expression of a particular filter dialect that evaluates to true or false and that
	// determines whether an instance of a CloudEvent will be delivered to the subscription's sink. If a filter
	// expression evaluates to false, the event MUST NOT be sent to the sink. If the expression evaluates to true, the
	// event MUST be attempted to be delivered. Support for particular filter dialects might vary across different
	// subscription managers. If a filter dialect is specified in a subscription that is unsupported by the subscription
	// manager, creation or updates of the subscription MUST be rejected with an error. See the Filter Dialects section
	// for further details.
	// +optional
	Filter *Filter `json:"filter,omitempty"`
}

type Protocol struct {
	// Protocol - Identifier of a delivery protocol. Because of WebSocket tunneling options for AMQP, MQTT and other
	// protocols, the URI scheme is not sufficient to identify the protocol. The protocols with existing CloudEvents
	// bindings are identified as "AMQP", "MQTT3", "MQTT5", "HTTP", "KAFKA", and "NATS". An implementation MAY add
	// support for further protocols.
	Protocol string `json:"protocol"`

	// ProtocolSettings - A set of settings specific to the selected delivery protocol provider. Options for those
	// settings are listed in the following subsection. An implementation MAY offer more options. Examples for such
	// settings are credentials, which generally vary by transport, rate limits and retry policies, or the QoS mode for
	// MQTT. See the Protocol Settings section for further details.
	// +optional
	ProtocolSettings *ProtocolSettings `json:"protocolsettings,omitempty"`
}

// Protocol holds the combined protocol settings.
type ProtocolSettings struct {
	Protocol       string `json:"-"`
	*HTTPProtocol  `json:",inline"`
	*MQTT3Protocol `json:",inline"`
	*MQTT5Protocol `json:",inline"`
	*AMQPProtocol  `json:",inline"`
	*KafkaProtocol `json:",inline"`
	*NATSProtocol  `json:",inline"`
}

func (p *Protocol) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(`{`)

	if pb, err := json.Marshal(p.Protocol); err != nil {
		return nil, err
	} else {
		b.WriteString(`"protocol":`)
		b.Write(pb)
	}

	if p.ProtocolSettings != nil {
		b.WriteString(`,"protocolsettings":`)
		switch p.Protocol {
		case "AMQP":
			if ps, err := json.Marshal(p.ProtocolSettings.AMQPProtocol); err != nil {
				return nil, err
			} else {
				b.Write(ps)
			}
		case "MQTT3":
			if ps, err := json.Marshal(p.ProtocolSettings.MQTT3Protocol); err != nil {
				return nil, err
			} else {
				b.Write(ps)
			}
		case "MQTT5":
			if ps, err := json.Marshal(p.ProtocolSettings.MQTT5Protocol); err != nil {
				return nil, err
			} else {
				b.Write(ps)
			}
		case "HTTP":
			if ps, err := json.Marshal(p.ProtocolSettings.HTTPProtocol); err != nil {
				return nil, err
			} else {
				b.Write(ps)
			}
		case "KAFKA":
			if ps, err := json.Marshal(p.ProtocolSettings.KafkaProtocol); err != nil {
				return nil, err
			} else {
				b.Write(ps)
			}
		case "NATS":
			if ps, err := json.Marshal(p.ProtocolSettings.NATSProtocol); err != nil {
				return nil, err
			} else {
				b.Write(ps)
			}
		default:
			return nil, fmt.Errorf("unsupported protocol: %q", p.Protocol)
		}
	}
	b.WriteString(`}`)

	return b.Bytes(), nil
}

// TODO I think we do not need a custom unmarshal. Test this.
//func (p *Protocol) UnmarshalJSON(b []byte) error {
//	var ref string
//	if err := json.Unmarshal(b, &ref); err != nil {
//		return err
//	}
//	r := ParseURIRef(ref)
//	if r != nil {
//		*u = *r
//	}
//	return nil
//}

type HTTPProtocol struct {
	// Headers - A set of key/value pairs that is copied into the HTTP request as custom headers.
	Headers map[string]string `json:"headers,omitempty"`

	// Method - The HTTP method to use for sending the message. This defaults to POST if not set.
	Method string `json:"method,omitempty"`
}

type MQTT3Protocol struct {
	// TopicName - The name of the MQTT topic to publish to.
	TopicName string `json:"topicname"`

	// QOS - MQTT quality of service (QoS) level: 0 (at most once), 1 (at least once), or 2 (exactly once).
	// This defaults to 1 if not set.
	// +optional
	QOS *int `json:"qos,omitempty"`

	// Retain - MQTT retain flag: true/false. This defaults to false if not set.
	// +optional
	Retain bool `json:"retain"`
}

type MQTT5Protocol struct {
	// TopicName - The name of the MQTT topic to publish to.
	TopicName string `json:"topicname"`

	// QOS - MQTT quality of service (QoS) level: 0 (at most once), 1 (at least once), or 2 (exactly once).
	// This defaults to 1 if not set.
	// +optional
	QOS *int `json:"qos,omitempty"`

	// Retain - MQTT retain flag: true/false. This defaults to false if not set.
	// +optional
	Retain bool `json:"retain"`

	// Expiry - MQTT expiry interval, in seconds. This value has no default value and the message will not expire if the
	// setting is absent.
	// +optional
	Expiry *int `json:"expiry,omitempty"`

	// UserProperties – A set of key/value pairs that are copied into the MQTT PUBLISH packet's user property section.
	// +optional
	UserProperties map[string]string `json:"userproperties,omitempty"`
}

type AMQPProtocol struct {
	// Address - The link target node in the AMQP container identified by the sink URI, if not expressed in the sink
	// URI's path portion.
	// +optional
	Address string `json:"address,omitempty"`

	// LinkName - Name to use for the AMQP link. If not set, a random link name is used.
	// +optional
	LinkName string `json:"linkname,omitempty"`

	// SenderSettlementMode - Allows to control the sender's settlement mode, which determines whether transfers are
	// performed "settled" (without acknowledgement) or "unsettled" (with acknowledgement). Default value is unsettled.
	// +optional
	SenderSettlementMode string `json:"sendersettlementmode,omitempty"`

	// LinkProperties - A set of key/value pairs that are copied into link properties for the send link.
	// +optional
	LinkProperties map[string]string `json:"linkproperties"`
}

type KafkaProtocol struct {
	// TopicName - The name of the Kafka topic to publish to.
	// +optional
	TopicName string `json:"topicname"`

	// PartitionKeyExtractor - A partition key extractor expression per the CloudEvents Kafka transport binding
	// specification.
	// +optional
	PartitionKeyExtractor string `json:"partitionkeyextractor,omitempty"`

	// ClientID
	// +optional
	ClientID string `json:"clientid,omitempty"`

	// ACKs
	// +optional
	ACKs string `json:"acks,omitempty"`
}

type NATSProtocol struct {
	// Subject - The name of the NATS subject to publish to.
	Subject string `json:"subject"`
}

type Filter struct {
	Dialect string        `json:"dialect"` // only "basic" is supported
	Filters []BasicFilter `json:"filters"`
}

// TODO: There could be other filters but they are currently not defined by the spec so we will do the easy way for now.
type BasicFilter struct {
	// Type - Value MUST be one of the following: prefix, suffix, exact.
	Type string `json:"type"`

	// The CloudEvents attribute (including extensions) to match the value indicated by the "value" property against.
	Property string `json:"property"`

	// Value - The value to match the CloudEvents attribute against. This expression is a string and matches are executed against the string representation of the attribute value.
	Value string `json:"value"`
}
