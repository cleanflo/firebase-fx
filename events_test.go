package register

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvents(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		assert.True(t, AuthenticationUserCreateEvent.Valid(), "Should be valid")
		assert.True(t, AuthenticationUserDeleteEvent.Valid(), "Should be valid")

		assert.True(t, FirestoreDocumentCreateEvent.Valid(), "Should be valid")
		assert.True(t, FirestoreDocumentDeleteEvent.Valid(), "Should be valid")
		assert.True(t, FirestoreDocumentUpdateEvent.Valid(), "Should be valid")
		assert.True(t, FirestoreDocumentWriteEvent.Valid(), "Should be valid")

		assert.True(t, PubSubPublishEvent.Valid(), "Should be valid")

		assert.True(t, RealtimeDBRefCreateEvent.Valid(), "Should be valid")
		assert.True(t, RealtimeDBRefDeleteEvent.Valid(), "Should be valid")
		assert.True(t, RealtimeDBRefUpdateEvent.Valid(), "Should be valid")
		assert.True(t, RealtimeDBRefWriteEvent.Valid(), "Should be valid")

		assert.True(t, StorageObjectArchiveEvent.Valid(), "Should be valid")
		assert.True(t, StorageObjectDeleteEvent.Valid(), "Should be valid")
		assert.True(t, StorageObjectFinalizeEvent.Valid(), "Should be valid")
		assert.True(t, StorageObjectMetadataUpdateEvent.Valid(), "Should be valid")

		assert.False(t, AuthEventType("").Valid(), "Should not be valid")
		assert.False(t, FirestoreEventType("").Valid(), "Should not be valid")
		assert.False(t, PubSubEventType("").Valid(), "Should not be valid")
		assert.False(t, RealtimeDBEventType("").Valid(), "Should not be valid")
		assert.False(t, StorageEventType("").Valid(), "Should not be valid")
	})
}
