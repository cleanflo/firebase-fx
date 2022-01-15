package register

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/functions/metadata"
	"github.com/stretchr/testify/assert"
)

type TestRTDBI struct {
	NNNNN float64
	BB    bool
	M     map[string]interface{}
	MS    struct {
		A string
		B int
	}
	G struct {
		Lat float64
		Lng float64
		S   struct {
			A string
			B int
		}
	}
	GNN struct {
		Lat float64
		Lng float64
	}
	T time.Time
}

func TestRTBD(t *testing.T) {
	reg := NewRegister()

	testCreatemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(RealtimeDBRefCreateEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/_/instances/[project-id]/refs/testColl/5914E2YLVWcUDHisQwQN",
		},
	})
	testUpdatemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(RealtimeDBRefUpdateEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/_/instances/[project-id]/refs/testColl/5914E2YLVWcUDHisQwQN",
		},
	})
	testDeletemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(RealtimeDBRefDeleteEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/_/instances/[project-id]/refs/testColl/5914E2YLVWcUDHisQwQN",
		},
	})
	testWritemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(RealtimeDBRefWriteEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/_/instances/[project-id]/refs/testColl/5914E2YLVWcUDHisQwQN",
		},
	})

	testDec := &Decoder{}
	err := json.Unmarshal([]byte(`{"data":{"NNNNN":112.45,"asd":"asd","hfghfd":false},"delta":{"NNNNN":112.46545}}`), testDec)
	if err != nil {
		t.Errorf("Error unmarshalling test rtdb data: %v", err)
	}

	testDbFunc := func(ctx context.Context, e RTDBEvent) error {
		assert.IsType(t, &TestRTDBI{}, e.Data, "Data did not match expected type")
		assert.IsType(t, &TestRTDBI{}, e.Delta, "Delta did not match expected type")

		assert.Contains(t, e.vars, "uid", "vars did not contain expected key")
		assert.Equal(t, "5914E2YLVWcUDHisQwQN", e.vars["uid"], "vars did not contain expected value")
		return nil
	}

	t.Run("Ref", func(t *testing.T) {
		db := reg.RealtimeDB()
		db.Ref("testRef")
		assert.Equalf(t, "testRef", db.resource, "Path should be testRef, got: %s", db.resource)

		db.Ref("sub")
		assert.Equalf(t, "testRef/sub", db.resource, "Path should be testRef/sub, got: %s", db.resource)

		db.Ref("deep/sub")
		assert.Equalf(t, "testRef/sub/deep/sub", db.resource, "Path should be testRef/sub/deep/sub, got: %s", db.resource)
	})

	t.Run("Register Create", func(t *testing.T) {
		db := reg.RealtimeDB().Ref("testColl/{uid}").Create(TestRTDBI{}, testDbFunc)
		assert.Same(t, db, reg.findRealtimeDB(RealtimeDBRefCreateEvent, "testColl/*"), "RealtimeDB function should be registered")
		assert.Same(t, db, reg.realtimeDB[RealtimeDBRefCreateEvent]["testColl/*"], "RealtimeDB function should be registered")
		assert.NotNil(t, db.fn, "Firestore function should be equal not nil")

		t.Log("RealtimeDB Function registered for CreateEvent")
	})

	t.Run("Register Delete", func(t *testing.T) {
		db := reg.RealtimeDB().Ref("testColl/{uid}").Delete(TestRTDBI{}, testDbFunc)
		assert.Same(t, db, reg.findRealtimeDB(RealtimeDBRefDeleteEvent, "testColl/*"), "RealtimeDB function should be registered")
		assert.Same(t, db, reg.realtimeDB[RealtimeDBRefDeleteEvent]["testColl/*"], "RealtimeDB function should be registered")
		assert.NotNil(t, db.fn, "Firestore function should be equal not nil")

		t.Log("RealtimeDB Function registered for CreateEvent")
	})

	t.Run("Register Update", func(t *testing.T) {
		db := reg.RealtimeDB().Ref("testColl/{uid}").Update(TestRTDBI{}, testDbFunc)
		assert.Same(t, db, reg.findRealtimeDB(RealtimeDBRefUpdateEvent, "testColl/*"), "RealtimeDB function should be registered")
		assert.Same(t, db, reg.realtimeDB[RealtimeDBRefUpdateEvent]["testColl/*"], "RealtimeDB function should be registered")
		assert.NotNil(t, db.fn, "Firestore function should be equal not nil")

		t.Log("RealtimeDB Function registered for CreateEvent")
	})

	t.Run("Register Write", func(t *testing.T) {
		db := reg.RealtimeDB().Ref("testColl/{uid}").Write(TestRTDBI{}, testDbFunc)
		assert.Same(t, db, reg.findRealtimeDB(RealtimeDBRefWriteEvent, "testColl/*"), "RealtimeDB function should be registered")
		assert.Same(t, db, reg.realtimeDB[RealtimeDBRefWriteEvent]["testColl/*"], "RealtimeDB function should be registered")
		assert.NotNil(t, db.fn, "RealtimeDB function should be equal not nil")

		t.Log("RealtimeDB Function registered for CreateEvent")
	})

	t.Run("Create Exec", func(t *testing.T) {
		cr := reg.RealtimeDB().Ref("testColl/{uid}").Create(TestRTDBI{}, testDbFunc)
		assert.NotNil(t, cr.fn, "RealtimeDB function should be not nil")

		err := cr.reg.EntryPoint(testCreatemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Delete Exec", func(t *testing.T) {
		de := reg.RealtimeDB().Ref("testColl/{uid}").Delete(TestRTDBI{}, testDbFunc)
		assert.NotNil(t, de.fn, "RealtimeDB function should be not nil")

		err := de.reg.EntryPoint(testDeletemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Update Exec", func(t *testing.T) {
		up := reg.RealtimeDB().Ref("testColl/{uid}").Update(TestRTDBI{}, testDbFunc)
		assert.NotNil(t, up.fn, "RealtimeDB function should be not nil")

		err := up.reg.EntryPoint(testUpdatemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Write Exec", func(t *testing.T) {
		wr := reg.RealtimeDB().Ref("testColl/{uid}").Write(TestRTDBI{}, testDbFunc)
		assert.NotNil(t, wr.fn, "RealtimeDB function should be not nil")

		err := wr.reg.EntryPoint(testWritemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

}
