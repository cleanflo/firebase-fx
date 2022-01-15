package register

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"cloud.google.com/go/functions/metadata"
)

// PubSub registers a function to the specified event, or returns the existing function if one already exists
func (f *FunctionRegistrar) PubSub(topic string) *PubSubFunction {
	if cf, ok := f.events[strings.ToLower(fmt.Sprintf("%s-%s", PubSubPublishEvent, topic))]; ok {
		if p, ok := cf.(*PubSubFunction); ok {
			return p
		}
	}

	p := &PubSubFunction{
		reg: f,
	}

	p.event = PubSubPublishEvent
	p.resource = topic
	p.reg.events[p.Name()] = p

	return p
}

// PubSubFunction is a wrapper for the PubSubFunc and the parent FunctionRegistrar
// Implements the CloudEventFunction interface
type PubSubFunction struct {
	cloudDeployer
	reg  *FunctionRegistrar
	fn   PubSubFunc
	data interface{}
}

// PubSubFunc is the function signature for the Pub/Sub CloudEvent
type PubSubFunc func(ctx context.Context, m PubSubMessage) error

// PubSubMessage is the expected payload for Pub/Sub CloudEvents.
type PubSubMessage struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}

// Publish registers the specified function to the specific topic for the Pub/Sub CloudEvent
// The provided data is used to populate the Data field of the PubSubMessage received by the function
//google.pubsub.topic.publish
func (p *PubSubFunction) Publish(data interface{}, fn PubSubFunc) *PubSubFunction {
	p.data = data
	p.fn = fn

	return p
}

// CloudEventFunction

// HandleCloudEvent handles the PubSub CloudEvent and calls the registered PubSubFunction
func (a *PubSubFunction) HandleCloudEvent(ctx context.Context, md *metadata.Metadata, dec *Decoder) error {
	msg := reflect.New(reflect.TypeOf(a.data)).Interface()

	err := dec.Decode(&msg)
	if err != nil {
		return Debug.Errf("failed to decode PubSubPublishEvent [%s]: %s: %s", md.EventType, err, string(dec.data))
	}

	m := PubSubMessage{
		Topic: a.resource,
		Data:  msg,
	}

	if a.fn != nil {
		err = a.fn(ctx, m)
		if err != nil {
			return Debug.Errf("registered PubSubFunc failed [%s]: %s: PubSubFunc %+v", md.EventType, err, a)
		}
	}
	return nil
}

// Name returns the name of the function: "pubsub.topic.publish/{topic}"
func (a *PubSubFunction) Name() string {
	return strings.ToLower(fmt.Sprintf("%s-%s", a.event, a.Resource()))
}

// Resource returns the resource of the function: "{topic}"
func (a *PubSubFunction) Resource() string {
	return a.resource
}

// Event returns the EventType of the function: PubSubPublishEvent
func (a *PubSubFunction) Event() EventType {
	return a.event.Type()
}
