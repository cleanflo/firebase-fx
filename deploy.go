package register

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/functions/metadata"
)

// CloudDeployFunction can be implemented by an event handler function
type CloudDeployFunction interface {
	HandleCloudEvent(context.Context, *metadata.Metadata, *Decoder) error
	Name() string
	Resource() string
	Event() EventType
}

type cloudDeployer struct {
	resource string
	event    event
}

func (f *FunctionRegistrar) Deploy() (s string) {
	projectId := f.projectID
	if projectId == "" {
		projectId = "$PROJECT_ID"
	}
	flags := func() (s string) {
		entrypoint := f.entrypoint
		if entrypoint == "" {
			entrypoint = "register.SharedEntrypoint"
		}
		s += fmt.Sprintf(" --entry-point %s", f.entrypoint)

		runtime := f.runtime
		if runtime == "" {
			runtime = "go116"
		}
		s += fmt.Sprintf(" --runtime %s", runtime)

		if projectId != "$PROJECT_ID" {
			s += fmt.Sprintf(" --project %s", projectId)
		}

		if f.verbosity != WarningVerbosity {
			s += fmt.Sprintf(" --verbosity %s", f.verbosity)
		}

		return s
	}

	// walk the functions and register each one
	cmds := []string{}
	for name, ev := range f.events {
		cmd := fmt.Sprintf("gcloud functions deploy %s \\\n", flags())
		switch ev.Event() {
		case AuthenticationUserCreateEvent.Type(), AuthenticationUserDeleteEvent.Type():
			cmd += "%s --trigger-event %s"
			cmds = append(cmds, fmt.Sprintf(cmd, name, ev.Event().String()))

		case FirestoreDocumentCreateEvent.Type(), FirestoreDocumentDeleteEvent.Type(), FirestoreDocumentUpdateEvent.Type(), FirestoreDocumentWriteEvent.Type():
			cmd += "%s --trigger-event %s --trigger-resource projects/%s/databases/(default)/documents/%s"
			cmds = append(cmds, fmt.Sprintf(cmd, name, ev.Event().String(), projectId, ev.Resource()))

		case PubSubPublishEvent.Type():
			cmd += "%s --trigger-topic %s"
			cmds = append(cmds, fmt.Sprintf(cmd, name, ev.Resource()))

		case RealtimeDBRefCreateEvent.Type(), RealtimeDBRefDeleteEvent.Type(), RealtimeDBRefUpdateEvent.Type(), RealtimeDBRefWriteEvent.Type():
			cmd += "%s --trigger-event %s --trigger-resource projects/_/instances/%s/refs/%s"
			cmds = append(cmds, fmt.Sprintf(cmd, name, ev.Event().String(), projectId, ev.Resource()))

		case StorageObjectFinalizeEvent.Type(), StorageObjectArchiveEvent.Type(), StorageObjectDeleteEvent.Type(), StorageObjectMetadataUpdateEvent.Type():
			cmd += "%s --trigger-event %s --trigger-resource %s"
			cmds = append(cmds, fmt.Sprintf(cmd, name, ev.Event().String(), ev.Resource()))
		}
	}

	// outputs a bash script the can be used to deploy the functions
	s = strings.Join(cmds, " &&  \\\n")
	return s
}

// VerbosityLevel is for setting the verbosity level for the deploy command
type VerbosityLevel int

func (v VerbosityLevel) String() string {
	switch v {
	case DebugVerbosity:
		return "debug"
	case InfoVerbosity:
		return "info"
	case WarningVerbosity:
		return "warning"
	case ErrorVerbosity:
		return "error"
	case CriticalVerbosity:
		return "critical"
	case NoneVerbosity:
		return "none"
	}

	return "warning"
}

const (
	DebugVerbosity VerbosityLevel = iota
	InfoVerbosity
	WarningVerbosity // default
	ErrorVerbosity
	CriticalVerbosity
	NoneVerbosity
)

// Verbosity sets the verbosity level for the deploy command
func (f *FunctionRegistrar) Verbosity(level VerbosityLevel) *FunctionRegistrar {
	f.verbosity = level
	return f
}

// WithProjectID sets the project id for the functions when deploying
// otherwise it will exclude the --project flag from the deploy command
// most often the project id is the simply the project name in kebab case
// eg. my-project-name
func (f *FunctionRegistrar) WithProjectID(id string) *FunctionRegistrar {
	f.projectID = id
	return f
}

// WithEntryPoint sets the entry point for the functions when deploying
// this is the actual name of the variable in your source code
// otherwise it will use the register.SharedEntryPoint
// again this assumes you have imported this package as register
// and you have registered your functions using register.Shared
func (f *FunctionRegistrar) WithEntrypoint(name string) *FunctionRegistrar {
	f.entrypoint = name
	return f
}

// Runtime is the runtime for the functions when deploying
type Runtime string

const (
	Go111 Runtime = "go111"
	Go113 Runtime = "go113"
	Go116 Runtime = "go116" // default
)

// WithRuntime sets the runtime for the functions when deploying
func (f *FunctionRegistrar) WithRuntime(runtime Runtime) *FunctionRegistrar {
	f.runtime = runtime
	return f
}

// FLAGS
// --entry-point ENTRY_POINT --project PROJECT_ID --verbosity VERBOSITY_LEVEL --runtime RUNTIME

// HTTP
// gcloud functions deploy FUNCTION_NAME --trigger-http --allow-unauthenticated

// PUBSUB
// gcloud functions deploy FUNCTION_NAME --trigger-topic TOPIC_NAME

// STORAGE
// gcloud functions deploy FUNCTION_NAME --trigger-event EVENT --trigger-resource YOUR_TRIGGER_BUCKET_NAME

// FIRESTORE
// gcloud functions deploy FUNCTION_NAME --trigger-event EVENT --trigger-resource "projects/PROJECT_ID/databases/(default)/documents/messages/{pushId}"

// REALTIME_DB
// gcloud functions deploy FUNCTION_NAME --trigger-event EVENT --trigger-resource projects/_/instances/PROJECT_ID/refs/messages/{pushId}/original

// AUTHENTICATION
// gcloud functions deploy FUNCTION_NAME --trigger-event EVENT

// ANALYTICS
// gcloud functions deploy FUNCTION_NAME --trigger-event EVENT  --trigger-resource projects/YOUR_PROJECT_ID/events/in_app_purchase

// REMOTE CONFIG
// gcloud functions deploy FUNCTION_NAME --trigger-event google.firebase.remoteconfig.update

// SCHEDULER
// gcloud pubsub topics create TOPIC_NAME
// gcloud functions deploy FUNCTION_NAME --trigger-topic TOPIC_NAME
// gcloud scheduler jobs create pubsub JOBNAME --topic TOPIC_NAME --schedule "every 5 minutes" --message CONFIRMATION_MESSAGE --time-zone TIME_ZONE
