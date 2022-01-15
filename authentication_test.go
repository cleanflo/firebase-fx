package register

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/functions/metadata"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	reg := NewRegister()

	testCreatemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(AuthenticationUserCreateEvent),
	})
	testDeletemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(AuthenticationUserDeleteEvent),
	})

	testDec := &Decoder{}
	err := json.Unmarshal([]byte(`{"email": "test@email.com", "metadata": {"createdAt": "2020-04-01T00:00:00Z"},"providerData": [{"email": "test@email.com", "providerId": "email", "uid": "test-uid" }],"uid": "test-uid"}`), testDec)
	if err != nil {
		t.Errorf("Error unmarshalling test authdata: %v", err)
	}

	testAuthFunc := func(ctx context.Context, e AuthEvent) error {
		assert.Equalf(t, "test@email.com", e.Email, "Email should match")
		assert.Equalf(t, "2020-04-01T00:00:00Z", e.Metadata.CreatedAt.Format(time.RFC3339), "CreatedAt should match")
		assert.Lenf(t, e.ProviderData, 1, "ProviderData should have 1 element, got: %d", len(e.ProviderData))
		assert.Equalf(t, "test@email.com", e.ProviderData[0].Email, "Email should match")
		assert.Equalf(t, "email", e.ProviderData[0].Provider, "Provider should match")
		assert.Equalf(t, "test-uid", e.ProviderData[0].UID, "UID should match")
		assert.Equalf(t, "test-uid", e.UID, "UID should match")
		return nil
	}

	t.Run("Register Create", func(t *testing.T) {
		auth := reg.Authentication().Create(testAuthFunc)
		assert.Same(t, auth, reg.events[auth.Name()], "Authentication function should be registered")
		assert.NotNil(t, auth.fn, "Authentication function should be equal not nil")

		t.Log("Authentication Function registered for UserCreateEvent")
	})

	t.Run("Register Delete", func(t *testing.T) {
		auth := reg.Authentication().Delete(testAuthFunc)
		assert.Same(t, auth, reg.events[auth.Name()], "Authentication function should be registered")
		assert.NotNil(t, auth.fn, "Authentication function should be equal not nil")

		t.Log("Authentication Function registered for UserDeleteEvent")
	})

	t.Run("Create Exec", func(t *testing.T) {
		create := reg.Authentication().Create(testAuthFunc)
		assert.NotNil(t, create.fn, "Authentication function should be not nil")

		err = create.reg.EntryPoint(testCreatemd, testDec)
		assert.Nil(t, err, "Authentication function error should be nil")
	})

	t.Run("Delete Exec", func(t *testing.T) {
		delete := reg.Authentication().Delete(testAuthFunc)
		assert.NotNil(t, delete.fn, "Authentication function should be not nil")

		err := delete.reg.EntryPoint(testDeletemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	testAuthErrFunc := func(ctx context.Context, e AuthEvent) error {
		return errors.New("test error")
	}

	t.Run("Authentication Exec Error", func(t *testing.T) {
		delete := reg.Authentication().Delete(testAuthErrFunc)
		assert.NotNil(t, delete.fn, "Authentication function should be not nil")

		err := delete.reg.EntryPoint(testDeletemd, testDec)
		assert.NotNil(t, err, "Error should be  not nil")

		create := reg.Authentication().Create(testAuthErrFunc)
		assert.NotNil(t, delete.fn, "Authentication function should be not nil")

		err = create.reg.EntryPoint(testDeletemd, testDec)
		assert.NotNil(t, err, "Error should be  not nil")
	})
}
