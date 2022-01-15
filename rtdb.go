package register

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"strings"

	"cloud.google.com/go/functions/metadata"
)

const (
	rtdbPathBase = "projects/_/instances/*/refs"
)

// RealtimeDB returns a new RealtimeDBFunction with the FunctionRegistrar set to the parent
func (f *FunctionRegistrar) RealtimeDB() *RealtimeDBFunction {
	s := &RealtimeDBFunction{reg: f}
	return s
}

// FindRealtimeDB attempts to match the event and path to a registered RealtimeDB function
// returns the function if found, otherwise nil
// expects the full path name as provided by the CloudEvent: "projects/_/instances/cleanflo-admin/refs/....."
func (f *FunctionRegistrar) findRealtimeDB(event RealtimeDBEventType, ref string) *RealtimeDBFunction {
	if f.realtimeDB[event] != nil && len(f.realtimeDB[event]) > 0 {
		ref = breakRef(ref)

		// collect the registered paths
		keys := make(pathKeys, 0, len(f.realtimeDB[event]))
		for k := range f.realtimeDB[event] {
			keys = append(keys, k)
		}

		if k := findPath(keys, ref); k != "" {
			fx := f.realtimeDB[event][k]
			return fx
		}
	}

	return nil
}

// Path creates a path for the RealtimeDB event that can be used for registering a function
// returns the path with all wildcard fields replaced with a *
// and saves a map of the segment positions var names for access within the function
func (r *RealtimeDBFunction) Path() string {
	m := make(map[int]string)

	pathParts := strings.Split(r.resource, "/")
	for i, part := range pathParts {
		if match := wildcardRegexp.MatchString(part); match {
			pathParts[i] = "*"
			m[i] = wildcard(part)
		} else if part == "*" {
			m[i] = part
		}
	}

	r.pathWildcards = m

	return path.Join(pathParts...)
}

// RealtimeDBFunction  is a wrapper for the expected data, RealtimeDBFunc and the parent FunctionRegistrar
// Implements the CloudEventFunction interface
type RealtimeDBFunction struct {
	cloudDeployer
	reg           *FunctionRegistrar
	fn            RealtimeDBFunc
	data          interface{}
	pathWildcards map[int]string
}

// RealtimeDBFunc is the function signature for firebase Realtime Database Cloud Events
type RealtimeDBFunc func(ctx context.Context, e RTDBEvent) error

// Ref appends the given path to the path that the RealtimeDBFunction is executed on
// Meant to be chained
func (r *RealtimeDBFunction) Ref(ref string) *RealtimeDBFunction {
	r.resource = path.Join(r.resource, ref)
	return r
}

// Write registers the specified function to the RefWriteEvent for firebase Realtime Database CloudEvent
// The provided data is used to populate .Data & .Delta of the RTDBEvent received by the function
//providers/google.firebase.database/eventTypes/ref.write
func (r *RealtimeDBFunction) Write(data interface{}, fn RealtimeDBFunc) *RealtimeDBFunction {
	r.fn = fn
	r.data = data

	if r.reg.realtimeDB[RealtimeDBRefWriteEvent] == nil {
		r.reg.realtimeDB[RealtimeDBRefWriteEvent] = make(map[string]*RealtimeDBFunction)
	}

	r.reg.realtimeDB[RealtimeDBRefWriteEvent][r.Path()] = r

	r.event = RealtimeDBRefWriteEvent
	r.reg.events[r.Name()] = r
	return r
}

// Create registers the specified function to the RefCreateEvent for firebase Realtime Database CloudEvent
// The provided data is used to populate .Data of the RTDBEvent received by the function
//providers/google.firebase.database/eventTypes/ref.create
func (r *RealtimeDBFunction) Create(data interface{}, fn RealtimeDBFunc) *RealtimeDBFunction {
	r.fn = fn
	r.data = data

	if r.reg.realtimeDB[RealtimeDBRefCreateEvent] == nil {
		r.reg.realtimeDB[RealtimeDBRefCreateEvent] = make(map[string]*RealtimeDBFunction)
	}

	r.reg.realtimeDB[RealtimeDBRefCreateEvent][r.Path()] = r

	r.event = RealtimeDBRefCreateEvent
	r.reg.events[r.Name()] = r
	return r
}

// Update registers the specified function to the RefUpdateEvent for firebase Realtime Database CloudEvent
// The provided data is used to populate .Data & .Delta of the RTDBEvent received by the function
//providers/google.firebase.database/eventTypes/ref.update
func (r *RealtimeDBFunction) Update(data interface{}, fn RealtimeDBFunc) *RealtimeDBFunction {
	r.fn = fn
	r.data = data

	if r.reg.realtimeDB[RealtimeDBRefUpdateEvent] == nil {
		r.reg.realtimeDB[RealtimeDBRefUpdateEvent] = make(map[string]*RealtimeDBFunction)
	}

	r.reg.realtimeDB[RealtimeDBRefUpdateEvent][r.Path()] = r

	r.event = RealtimeDBRefUpdateEvent
	r.reg.events[r.Name()] = r
	return r
}

// Delete registers the specified function to the RefDeleteEvent for firebase Realtime Database CloudEvent
// The provided data is used to populate .Delta of the RTDBEvent received by the function
//providers/google.firebase.database/eventTypes/ref.delete
func (r *RealtimeDBFunction) Delete(data interface{}, fn RealtimeDBFunc) *RealtimeDBFunction {
	r.fn = fn
	r.data = data

	if r.reg.realtimeDB[RealtimeDBRefDeleteEvent] == nil {
		r.reg.realtimeDB[RealtimeDBRefDeleteEvent] = make(map[string]*RealtimeDBFunction)
	}

	r.reg.realtimeDB[RealtimeDBRefDeleteEvent][r.Path()] = r

	r.event = RealtimeDBRefDeleteEvent
	r.reg.events[r.Name()] = r
	return r
}

// RTDBEvent is the expected payload of a firebase Realtime Database CloudEvent
type RTDBEvent struct {
	vars  map[string]string
	Data  interface{}
	Delta interface{}
}

// Vars returns the map of variable name that were registered on the path that matched this function
func (e *RTDBEvent) Vars() map[string]string {
	return e.vars
}

// CloudEventFunction

// HandleCloudEvent handles the Firebase RealtimeDB CloudEvent and calls the registered RealtimeDBFunction
func (a *RealtimeDBFunction) HandleCloudEvent(ctx context.Context, md *metadata.Metadata, dec *Decoder) error {
	evt := RTDBEvent{}

	if a == nil {
		return Debug.Errf("no RealtimeDBFunc registered for [%s]: %s", md.EventType, md.Resource.RawPath)
	}

	var reqData struct {
		Data  json.RawMessage `json:"data"`
		Delta json.RawMessage `json:"delta"`
	}

	err := dec.Decode(&reqData)
	if err != nil {
		return Debug.Errf("failed to decode realtimeDB event [%s]: %s: %s", md.EventType, err, string(dec.data))
	}

	dataType := reflect.TypeOf(a.data)
	dataT := reflect.New(dataType)
	deltaT := reflect.New(dataType)

	err = json.Unmarshal(reqData.Data, dataT.Interface())
	if err != nil {
		return Debug.Errf("failed to unmarshal realtimeDB event data [%s]: %s: %s", md.EventType, err, string(dec.data))
	}

	err = json.Unmarshal(reqData.Delta, deltaT.Interface())
	if err != nil {
		return Debug.Errf("failed to unmarshal realtimeDB event delta [%s]: %s: %s", md.EventType, err, string(dec.data))
	}

	evt.Data = dataT.Interface()
	evt.Delta = deltaT.Interface()

	evt.vars = extractVars(breakRef(md.Resource.RawPath), a.pathWildcards)

	err = a.fn(ctx, evt)
	if err != nil {
		return Debug.Errf("registered realtimeDBFunc failed [%s]: %s: RealtimeDBFunc %+v", md.EventType, err, a)
	}
	return nil
}

// Name returns the name of the function: "rtdb.ref.{write,create,delete,update}-{path-segments}"
func (a *RealtimeDBFunction) Name() string {
	return fmt.Sprintf("%s-%s", a.event, normalizeRegexp.ReplaceAllString(a.Resource(), "-$1"))
}

// Resource returns the resource of the function: "ref/{ref-id}/...."
func (a *RealtimeDBFunction) Resource() string {
	return a.resource
}

// Event returns the EventType of the function:
// RealtimeDBRefWriteEvent / RealtimeDBRefCreateEvent / RealtimeDBRefUpdateEvent / RealtimeDBRefDeleteEvent
func (a *RealtimeDBFunction) Event() EventType {
	return a.event.Type()
}
