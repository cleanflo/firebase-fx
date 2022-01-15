package register

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/functions/metadata"
	"github.com/stretchr/testify/assert"
)

type TestPubSubI struct {
	Email string `json:"email"`
}

func TestPubSub(t *testing.T) {
	reg := NewRegister()
	testmd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(PubSubPublishEvent),
	})

	testDec := &Decoder{}
	err := json.Unmarshal([]byte(`{"topic": "test-topic", "email": "other@email.com"}`), testDec)
	if err != nil {
		t.Errorf("Error unmarshalling test pubsub data: %v", err)
	}

	testPubSubFunc := func(ctx context.Context, e PubSubMessage) error {
		assert.IsType(t, &TestPubSubI{}, e.Data, "Data should be of type TestPubSubI")

		if v, ok := e.Data.(*TestPubSubI); ok {
			assert.Equalf(t, "other@email.com", v.Email, "Email should match")
		} else {
			t.Errorf("Data should be of type TestPubSub, got: %T", e.Data)
		}

		assert.Equalf(t, "test-topic", e.Topic, "Email should match")
		return nil
	}

	t.Run("Test Register Publish", func(t *testing.T) {
		ps := reg.PubSub("test-topic").Publish(TestPubSubI{}, testPubSubFunc)

		assert.Same(t, ps, reg.PubSub("test-topic"), "PubSub function should be registered")
		assert.Same(t, ps, reg.events[strings.ToLower(fmt.Sprintf("%s-%s", PubSubPublishEvent, "test-topic"))], "PubSub function should be registered")
		assert.NotNil(t, ps.fn, "PubSub function should be equal not nil")

		t.Log("PubSub Function registered for UserCreateEvent")
	})

	t.Run("PubSub Exec", func(t *testing.T) {
		ps := reg.PubSub("test-topic").Publish(TestPubSubI{}, testPubSubFunc)
		assert.NotNil(t, ps.fn, "PubSub function should be not nil")
		assert.NotNil(t, ps.data, "PubSub data should be not nil")

		err := ps.reg.EntryPoint(testmd, testDec)
		assert.Nil(t, err, "PubSub function error should be nil")
	})

	errmd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(PubSubPublishEvent),
	})

	testPubSubErrFunc := func(ctx context.Context, e PubSubMessage) error {
		return errors.New("test error")
	}

	testerrDec := &Decoder{}
	err = json.Unmarshal([]byte(`{"topic": "non-topic", "email": "other@email.com"}`), testerrDec)
	if err != nil {
		t.Errorf("Error unmarshalling test pubsub data: %v", err)
	}

	t.Run("PubSub Exec Error", func(t *testing.T) {
		ps := reg.PubSub("non-topic").Publish(TestPubSubI{}, testPubSubErrFunc)
		assert.NotNil(t, ps.fn, "PubSub function should be not nil")
		assert.NotNil(t, ps.data, "PubSub data should be not nil")

		err = ps.reg.EntryPoint(errmd, testerrDec)
		assert.NotNilf(t, err, "PubSub function error should be nil: %s", err)
	})

}
