package register

type event interface {
	Type() EventType
	String() string
	Valid() bool
}

type EventType string

// String returns the event type
func (e EventType) String() string {
	return string(e)
}

const (
	// Google Analytics Firebase event types
	AnalyticsLogEvent AnalyticsEventType = "providers/google.firebase.analytics/eventTypes/event.log"

	// Authentication event types
	AuthenticationUserCreateEvent AuthEventType = "providers/firebase.auth/eventTypes/user.create"
	AuthenticationUserDeleteEvent AuthEventType = "providers/firebase.auth/eventTypes/user.delete"

	// Firestore event types
	FirestoreDocumentCreateEvent FirestoreEventType = "providers/cloud.firestore/eventTypes/document.create"
	FirestoreDocumentDeleteEvent FirestoreEventType = "providers/cloud.firestore/eventTypes/document.delete"
	FirestoreDocumentUpdateEvent FirestoreEventType = "providers/cloud.firestore/eventTypes/document.update"
	FirestoreDocumentWriteEvent  FirestoreEventType = "providers/cloud.firestore/eventTypes/document.write"

	// Pub/Sub event types
	PubSubPublishEvent PubSubEventType = "google.pubsub.topic.publish"

	// Realtime Database event types
	RealtimeDBRefCreateEvent RealtimeDBEventType = "providers/google.firebase.database/eventTypes/ref.create"
	RealtimeDBRefDeleteEvent RealtimeDBEventType = "providers/google.firebase.database/eventTypes/ref.delete"
	RealtimeDBRefUpdateEvent RealtimeDBEventType = "providers/google.firebase.database/eventTypes/ref.update"
	RealtimeDBRefWriteEvent  RealtimeDBEventType = "providers/google.firebase.database/eventTypes/ref.write"

	// Firebase Remote Config event types
	RemoteConfigUpdateEvent RemoteConfigEventType = "remoteConfig.update"

	// Scheduler event types
	SchedulerRunEvent SchedulerEventType = "google.pubsub.topic.publish"

	// Storage event types
	StorageObjectArchiveEvent        StorageEventType = "google.storage.object.archive"
	StorageObjectDeleteEvent         StorageEventType = "google.storage.object.delete"
	StorageObjectFinalizeEvent       StorageEventType = "google.storage.object.finalize"
	StorageObjectMetadataUpdateEvent StorageEventType = "google.storage.object.metadataUpdate"
)

// AnalyticsEventType is the event type for Analytics CloudEvents
type AnalyticsEventType EventType

// Valid reports whether the AnalyticsEventType is valid
func (a AnalyticsEventType) Valid() bool {
	return a == AnalyticsLogEvent
}

// Type returns the event type
func (e AnalyticsEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e AnalyticsEventType) String() string {
	return "analyticsLog"
}

// AuthEventType is the event type for Firebase Authentication CloudEvents
type AuthEventType EventType

// Valid reports whether the AuthEventType is valid
func (a AuthEventType) Valid() bool {
	switch a {
	case AuthenticationUserCreateEvent, AuthenticationUserDeleteEvent:
		return true
	}
	return false
}

// Type returns the event type
func (e AuthEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e AuthEventType) String() string {
	switch e {
	case AuthenticationUserCreateEvent:
		return "authUserCreate"
	case AuthenticationUserDeleteEvent:
		return "authUserDelete"
	}
	return "authUserUunknown"
}

// FirestoreEventType is the event type for Firestore CloudEvents
type FirestoreEventType EventType

// Valid reports whether the FirstoreEventType is valid
func (f FirestoreEventType) Valid() bool {
	switch f {
	case FirestoreDocumentCreateEvent, FirestoreDocumentUpdateEvent, FirestoreDocumentDeleteEvent, FirestoreDocumentWriteEvent:
		return true
	}
	return false
}

// Type returns the event type
func (e FirestoreEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e FirestoreEventType) String() string {
	switch e {
	case FirestoreDocumentCreateEvent:
		return "firestoreDocCreate"
	case FirestoreDocumentDeleteEvent:
		return "firestoreDocDelete"
	case FirestoreDocumentUpdateEvent:
		return "firestoreDocUpdate"
	case FirestoreDocumentWriteEvent:
		return "firestoreDocWrite"
	}
	return "firestoreDocUnknown"
}

// PubSubEventType is the event type for Pub/Sub CloudEvents
type PubSubEventType EventType

// Valid reports whether the PubSubEventType is valid
func (p PubSubEventType) Valid() bool {
	return p == PubSubPublishEvent
}

// Type returns the event type
func (e PubSubEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e PubSubEventType) String() string {
	return "pubsubPublish"
}

// RealtimeDBEventType is the event type for Firebase Realtime Database CloudEvents
type RealtimeDBEventType EventType

// Valid reports whether the RealtimeDBEventType is valid
func (r RealtimeDBEventType) Valid() bool {
	switch r {
	case RealtimeDBRefWriteEvent, RealtimeDBRefCreateEvent, RealtimeDBRefUpdateEvent, RealtimeDBRefDeleteEvent:
		return true
	}
	return false
}

// Type returns the event type
func (e RealtimeDBEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e RealtimeDBEventType) String() string {
	switch e {
	case RealtimeDBRefCreateEvent:
		return "rtdbRefCreate"
	case RealtimeDBRefDeleteEvent:
		return "rtdbRefDelete"
	case RealtimeDBRefUpdateEvent:
		return "rtdbRefUpdate"
	case RealtimeDBRefWriteEvent:
		return "rtdbRefWrite"
	}
	return "rtdbRefUnknown"
}

// RemoteConfigEventType is the event type for Firebase Remote Config CloudEvents
type RemoteConfigEventType EventType

// Valid reports whether the RemoteConfigEventType is valid
func (r RemoteConfigEventType) Valid() bool {
	return r == RemoteConfigUpdateEvent
}

// Type returns the event type
func (e RemoteConfigEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e RemoteConfigEventType) String() string {
	return "remoteConfigUpdate"
}

// StorageEventType is the event type for Storage CloudEvents
type StorageEventType EventType

// Valid reports whether the StorageEventType is valid
func (s StorageEventType) Valid() bool {
	switch s {
	case StorageObjectFinalizeEvent, StorageObjectDeleteEvent, StorageObjectArchiveEvent, StorageObjectMetadataUpdateEvent:
		return true
	}
	return false
}

// Type returns the event type
func (e StorageEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e StorageEventType) String() string {
	switch e {
	case StorageObjectArchiveEvent:
		return "storageObjectArchive"
	case StorageObjectDeleteEvent:
		return "storageObjectDelete"
	case StorageObjectFinalizeEvent:
		return "storageObjectFinalize"
	case StorageObjectMetadataUpdateEvent:
		return "storageObjectMetadata"
	}
	return "storageObjectUnknown"
}

type SchedulerEventType PubSubEventType

// Valid reports whether the RemoteConfigEventType is valid
func (r SchedulerEventType) Valid() bool {
	return r == SchedulerEventType(PubSubPublishEvent)
}

// Type returns the event type
func (e SchedulerEventType) Type() EventType {
	return EventType(e)
}

// String returns a minimal/readable representation of the event type
func (e SchedulerEventType) String() string {
	return "schedulerRun"
}
