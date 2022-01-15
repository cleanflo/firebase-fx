package register

import (
	"context"
	"time"

	"cloud.google.com/go/functions/metadata"
)

// Authentication returns a new AuthenticationFunction with the FunctionRegistrar set to the parent
func (f *FunctionRegistrar) Authentication() *AuthenticationFunction {
	a := &AuthenticationFunction{reg: f}
	return a
}

// AuthenticationFunction is a wrapper for the AuthenticationFunc and the parent FunctionRegistrar
// Implements the CloudEventFunction interface
type AuthenticationFunction struct {
	cloudDeployer
	reg *FunctionRegistrar
	fn  AuthenticationFunc
}

// AuthenticationFunc is the function signature for the Firebase Authentication CloudEvent
type AuthenticationFunc func(ctx context.Context, e AuthEvent) error

// AuthEvent is the expected payload for Firestore Auth CloudEvents.
type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	ProviderData []struct {
		Email    string `json:"email"`
		Provider string `json:"providerId"`
		UID      string `json:"uid"`
	} `json:"providerData"`
	UID string `json:"uid"`
}

// Create registers the specified function to the UserCreated event for the Firebase Authentication CloudEvent
//providers/firebase.auth/eventTypes/user.create
func (a *AuthenticationFunction) Create(fn AuthenticationFunc) *AuthenticationFunction {
	a.fn = fn

	a.event = AuthenticationUserCreateEvent
	a.reg.events[a.Name()] = a
	return a
}

// Delete registers the specified function to the UserDeleted event for the Firebase Authentication CloudEvent
//providers/firebase.auth/eventTypes/user.delete
func (a *AuthenticationFunction) Delete(fn AuthenticationFunc) *AuthenticationFunction {
	a.fn = fn

	a.event = AuthenticationUserDeleteEvent
	a.reg.events[a.Name()] = a
	return a
}

// CloudEventFunction

// HandleCloudEvent handles the Firebase Authentication CloudEvent and calls the registered AuthenticationFunction
func (a *AuthenticationFunction) HandleCloudEvent(ctx context.Context, md *metadata.Metadata, dec *Decoder) error {
	event := &AuthEvent{}
	err := dec.Decode(&event)
	if err != nil {
		return Debug.Errf("failed to decode AuthEvent [%s]: %s: %s", md.EventType, err, string(dec.data))
	}

	if a.fn != nil {
		err = a.fn(ctx, *event)
		if err != nil {
			return Debug.Errf("registered AuthFunc failed [%s]: %s: AuthFunc %+v", md.EventType, err, a)
		}
	}

	return nil
}

// Name returns the name of the function: "auth.user.{create,delete}"
func (a *AuthenticationFunction) Name() string {
	return a.event.String()
}

// Resource returns the resource of the function: "providers/firebase.auth/eventTypes/user.{create,delete}"
func (a *AuthenticationFunction) Resource() string {
	return a.Event().String()
}

// Event returns the EventType of the function: AuthenticationUserCreateEvent / AuthenticationUserDeleteEvent
func (a *AuthenticationFunction) Event() EventType {
	return a.event.Type()
}
