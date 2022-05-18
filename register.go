package register

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/functions/metadata"
	"github.com/gorilla/mux"
)

var (
	Shared               *FunctionRegistrar = NewRegister()         // the shared registrar
	SharedEntryPoint                        = Shared.EntryPoint     // the entrypoint for shared background functions
	SharedHttpEntrypoint                    = Shared.HttpEntrypoint // the entrypoint for shared http functions
)

// Decoder is a helper to retrieve the JSON data from the request
type Decoder struct {
	data []byte
}

// UnmarshalJSON is a custom marshaller for the Decoder
func (d *Decoder) UnmarshalJSON(b []byte) error {
	d.data = b
	return nil
}

// Decode decodes the event data into the given type
func (d *Decoder) Decode(v interface{}) error {
	return json.Unmarshal(d.data, v)
}

// FunctionRegistrar is the registrar for functions
type FunctionRegistrar struct {
	http     *mux.Router // mapped by route
	handlers map[string]*HttpFunction

	// authentication map[AuthEventType]*AuthenticationFunction              // mapped by event type
	firestore  map[FirestoreEventType]map[string]*FirestoreFunction   // mapped by event type & path
	realtimeDB map[RealtimeDBEventType]map[string]*RealtimeDBFunction // mapped by event type & path
	// pubsub         map[string]*PubSubFunction                             // mapped by topic
	storage map[StorageEventType]map[string]*StorageFunction // mapped by event type & path

	// TODO:
	// analytics    map[AnalyticsEventType]*AnalyticsFunction       // mapped by event type
	// remoteConfig map[RemoteConfigEventType]*RemoteConfigFunction // mapped by event type
	// scheduler    map[string]*SchedulerFunction                   // mapped by event type

	events    map[string]CloudDeployFunction // pointers to all functions with a generated name
	projectID string
	registrar string
	verbosity VerbosityLevel
	runtime   Runtime

	httpUnauthenticated bool
}

// NewRegister creates a new registrar with all top level maps initialized
// nested maps are intialized when a function is registered
func NewRegister() *FunctionRegistrar {
	return &FunctionRegistrar{
		http:     mux.NewRouter(),
		handlers: make(map[string]*HttpFunction),
		events:   make(map[string]CloudDeployFunction),
		// pubsub:         make(map[string]*PubSubFunction),
		storage:    make(map[StorageEventType]map[string]*StorageFunction),
		firestore:  make(map[FirestoreEventType]map[string]*FirestoreFunction),
		realtimeDB: make(map[RealtimeDBEventType]map[string]*RealtimeDBFunction),
		// authentication: make(map[AuthEventType]*AuthenticationFunction),
	}
}

// The entrypoint for all functions.
// This is the function that is called when the function is triggered.
//     	 	f := register.NewRegister()
//          f.PubSub("topic").Publish(&mydata, myPubsubFunc)
// now f.Entrypoint can be used as the entrypoint when deploying the function
func (f *FunctionRegistrar) EntryPoint(ctx context.Context, dec *Decoder) error {
	md, err := metadata.FromContext(ctx)
	if err != nil {
		return Debug.Err("context metadata failed: %s", err)
	}

	switch true {

	case FirestoreEventType(md.EventType).Valid():
		fsFunc := f.findFirestore(FirestoreEventType(md.EventType), md.Resource.RawPath)
		if fsFunc == nil {
			return Debug.Errf("no FirestoreFunc registered for [%s]: %s", md.EventType, md.Resource.RawPath)
		}

		fsFunc.HandleCloudEvent(ctx, md, dec)

		return nil

	case AuthEventType(md.EventType).Valid():
		if c, ok := f.findEvent(AuthEventType(md.EventType).String()); ok {
			err = c.HandleCloudEvent(ctx, md, dec)
			if err != nil {
				return Debug.Errf("registered authFunc failed [%s]: %s: AuthFunc %+v", md.EventType, err, c)
			}
		}

		return nil
	case PubSubEventType(md.EventType).Valid():
		var m PubSubMessage
		err = dec.Decode(&m)
		if err != nil {
			return Debug.Errf("failed to decode topic [%s]: %s: %s", md.EventType, err, string(dec.data))
		}

		pubFunc := f.PubSub(m.Topic)
		if err := pubFunc.HandleCloudEvent(ctx, md, dec); err != nil {
			return Debug.Err("failed to handle cloud event", err)
		}

		return nil

	case RealtimeDBEventType(md.EventType).Valid():
		dbFunc := f.findRealtimeDB(RealtimeDBEventType(md.EventType), md.Resource.RawPath)
		if err := dbFunc.HandleCloudEvent(ctx, md, dec); err != nil {
			return Debug.Err("failed to handle cloud event", err)
		}
		return nil

	case StorageEventType(md.EventType).Valid():
		stFunc := f.findStorage(StorageEventType(md.EventType), md.Resource.Name)
		if err := stFunc.HandleCloudEvent(ctx, md, dec); err != nil {
			return Debug.Err("failed to handle cloud event", err)
		}

		return nil
	}

	return Debug.Errf("registrar closed without any results: event type did not decode: %s: %s", md.EventType, string(dec.data))
}

// Find locates a registered function by name.
func (f *FunctionRegistrar) findEvent(name string) (CloudDeployFunction, bool) {
	if v, ok := f.events[name]; ok {
		return v, ok
	}

	return nil, false
}
