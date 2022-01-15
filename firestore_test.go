package register

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/functions/metadata"
	"github.com/stretchr/testify/assert"
)

/*
{"oldValue":{"createTime":"2022-01-02T22:19:55.897215Z","fields":{"BB":{"booleanValue":true},"M":{"mapValue":{"fields":{"Asd":{"integerValue":"2355"},"Bfg":{"booleanValue":true}}}},"MS":{"mapValue":{"fields":{"A":{"stringValue":"ertert"},"B":{"integerValue":"6544"}}}},"NNNNN":{"integerValue":"123"},"geo":{"geoPointValue":{"latitude":50.55,"longitude":-104.87}}},"name":"projects/cleanflo-admin/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN","updateTime":"2022-01-03T03:39:46.371407Z"},"updateMask":{"fieldPaths":["MS.B"]},"value":{"createTime":"2022-01-02T22:19:55.897215Z","fields":{"BB":{"booleanValue":true},"M":{"mapValue":{"fields":{"Asd":{"integerValue":"2355"},"Bfg":{"booleanValue":true}}}},"MS":{"mapValue":{"fields":{"A":{"stringValue":"ertert"},"B":{"integerValue":"1234"}}}},"NNNNN":{"integerValue":"123"},"geo":{"geoPointValue":{"latitude":50.55,"longitude":-104.87}}},"name":"projects/cleanflo-admin/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN","updateTime":"2022-01-03T03:41:22.930655Z"}}

{
	"oldValue":{
		"createTime":"2022-01-02T22:19:55.897215Z",
		"fields":{
			"BB":{"booleanValue":true},
			"M":{
				"mapValue":{"fields":{
					"Asd":{"integerValue":"2355"},
					"Bfg":{"booleanValue":true}
				}}
			},
			"MS":{
				"mapValue":{"fields":{
					"A":{"stringValue":"ertert"},
					"B":{"integerValue":"6544"}
				}}
			},
			"NNNNN":{"integerValue":"123"},
			"geo":{"geoPointValue":{"latitude":50.55,"longitude":-104.87}}
		},
		"name":"projects/cleanflo-admin/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN",
		"updateTime":"2022-01-03T03:39:46.371407Z"
	},
	"updateMask":{
		"fieldPaths":["MS.B"]
	},
	"value":{
		"createTime":"2022-01-02T22:19:55.897215Z",
		"fields":{
			"BB":{"booleanValue":true},
			"M":{
				"mapValue":{"fields":{
					"Asd":{"integerValue":"2355"},
					"Bfg":{"booleanValue":true}
				}}
			},
			"MS":{
				"mapValue":{"fields":{
					"A":{"stringValue":"ertert"},
					"B":{"integerValue":"1234"}
				}}
			},
			"NNNNN":{"integerValue":"123"},
			"geo":{"geoPointValue":{"latitude":50.55,"longitude":-104.87}}
		},
		"name":"projects/cleanflo-admin/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN",
		"updateTime":"2022-01-03T03:41:22.930655Z"
	}
}
*/

type TestFirestoreI struct {
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

func TestFirestore(t *testing.T) {
	reg := NewRegister()

	testCreatemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(FirestoreDocumentCreateEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/[project-name]/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN",
		},
	})
	testDeletemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(FirestoreDocumentDeleteEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/[project-name]/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN",
		},
	})
	testUpdatemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(FirestoreDocumentUpdateEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/[project-name]/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN",
		},
	})
	testWritemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(FirestoreDocumentWriteEvent),
		Resource: &metadata.Resource{
			RawPath: "projects/[project-name]/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN",
		},
	})

	testDec := &Decoder{}
	err := json.Unmarshal([]byte(`{"oldValue":{"createTime":"2022-01-02T22:19:55.897215Z","fields":{"BB":{"booleanValue":true},"M":{"mapValue":{"fields":{"Asd":{"integerValue":"2355"},"Bfg":{"booleanValue":true}}}},"MS":{"mapValue":{"fields":{"A":{"stringValue":"ertert"},"B":{"integerValue":"6544"}}}},"NNNNN":{"integerValue":"123"},"geo":{"geoPointValue":{"latitude":50.55,"longitude":-104.87}}},"name":"projects/cleanflo-admin/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN","updateTime":"2022-01-03T03:39:46.371407Z"},"updateMask":{"fieldPaths":["MS.B"]},"value":{"createTime":"2022-01-02T22:19:55.897215Z","fields":{"BB":{"booleanValue":true},"M":{"mapValue":{"fields":{"Asd":{"integerValue":"2355"},"Bfg":{"booleanValue":true}}}},"MS":{"mapValue":{"fields":{"A":{"stringValue":"ertert"},"B":{"integerValue":"1234"}}}},"NNNNN":{"integerValue":"123"},"geo":{"geoPointValue":{"latitude":50.55,"longitude":-104.87}}},"name":"projects/cleanflo-admin/databases/(default)/documents/testColl/5914E2YLVWcUDHisQwQN","updateTime":"2022-01-03T03:41:22.930655Z"}}`), testDec)
	if err != nil {
		t.Errorf("Error unmarshalling test firestore data: %v", err)
	}

	testFsFunc := func(ctx context.Context, e FirestoreEvent) error {
		assert.IsType(t, &TestFirestoreI{}, e.OldValue.Fields, "OldValue.Fields did not match expected type")
		assert.IsType(t, &TestFirestoreI{}, e.Value.Fields, "Value.Fields did not match expected type")

		assert.Contains(t, e.vars, "uid", "vars did not contain expected key")
		assert.Equal(t, "5914E2YLVWcUDHisQwQN", e.vars["uid"], "vars did not contain expected value")
		return nil
	}

	t.Run("Collection/Document", func(t *testing.T) {
		fs := reg.Firestore()
		fs.Collection("testColl").Document("testDoc")
		assert.Equalf(t, "testColl/testDoc", fs.resource, "Path should be testColl/testDoc, got: %s", fs.resource)

		fs.Collection("subColl")
		assert.Equalf(t, "testColl/testDoc/subColl", fs.resource, "Path should be testColl/testDoc/subColl, got: %s", fs.resource)

		fs.Collection("nonColl")
		assert.Equalf(t, "testColl/testDoc/nonColl", fs.resource, "Path should be testColl/testDoc/nonColl, got: %s", fs.resource)

		fs.Document("nonDoc")
		assert.Equalf(t, "testColl/testDoc/nonColl/nonDoc", fs.resource, "Path should be testColl/testDoc/nonColl/nonDoc, got: %s", fs.resource)

		fs.Document("subDoc")
		assert.Equalf(t, "testColl/testDoc/nonColl/subDoc", fs.resource, "Path should be testColl/testDoc/nonColl/subDoc, got: %s", fs.resource)
	})

	t.Run("Register Create", func(t *testing.T) {
		fs := reg.Firestore().Collection("testColl").Document("{uid}").Create(TestFirestoreI{}, testFsFunc)
		assert.Same(t, fs, reg.findFirestore(FirestoreDocumentCreateEvent, "testColl/*"), "Firestore function should be registered")
		assert.Same(t, fs, reg.firestore[FirestoreDocumentCreateEvent]["testColl/*"], "Firestore function should be registered")
		assert.NotNil(t, fs.fn, "Firestore function should be equal not nil")

		t.Log("Firestore Function registered for CreateEvent")
	})

	t.Run("Register Delete", func(t *testing.T) {
		fs := reg.Firestore().Collection("testColl").Document("{uid}").Delete(TestFirestoreI{}, testFsFunc)
		assert.Same(t, fs, reg.findFirestore(FirestoreDocumentDeleteEvent, "testColl/*"), "Firestore function should be registered")
		assert.Same(t, fs, reg.firestore[FirestoreDocumentDeleteEvent]["testColl/*"], "Firestore function should be registered")
		assert.NotNil(t, fs.fn, "Firestore function should be equal not nil")

		t.Log("Firestore Function registered for DeleteEvent")
	})

	t.Run("Register Update", func(t *testing.T) {
		fs := reg.Firestore().Collection("testColl").Document("{uid}").Update(TestFirestoreI{}, testFsFunc)
		assert.Same(t, fs, reg.findFirestore(FirestoreDocumentUpdateEvent, "testColl/*"), "Firestore function should be registered")
		assert.Same(t, fs, reg.firestore[FirestoreDocumentUpdateEvent]["testColl/*"], "Firestore function should be registered")
		assert.NotNil(t, fs.fn, "Firestore function should be equal not nil")

		t.Log("Firestore Function registered for UpdateEvent")
	})

	t.Run("Register Write", func(t *testing.T) {
		fs := reg.Firestore().Collection("testColl").Document("{uid}").Write(TestFirestoreI{}, testFsFunc)
		assert.Same(t, fs, reg.findFirestore(FirestoreDocumentWriteEvent, "testColl/*"), "Firestore function should be registered")
		assert.Same(t, fs, reg.firestore[FirestoreDocumentWriteEvent]["testColl/*"], "Firestore function should be registered")
		assert.NotNil(t, fs.fn, "Firestore function should be equal not nil")

		t.Log("Firestore Function registered for WriteEvent")
	})

	t.Run("Create Exec", func(t *testing.T) {
		cr := reg.Firestore().Collection("testColl").Document("{uid}").Create(TestFirestoreI{}, testFsFunc)
		assert.NotNil(t, cr.fn, "Firestore function should be not nil")

		err := cr.reg.EntryPoint(testCreatemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Delete Exec", func(t *testing.T) {
		de := reg.Firestore().Collection("testColl").Document("{uid}").Delete(TestFirestoreI{}, testFsFunc)
		assert.NotNil(t, de.fn, "Firestore function should be not nil")

		err := de.reg.EntryPoint(testDeletemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Update Exec", func(t *testing.T) {
		up := reg.Firestore().Collection("testColl").Document("{uid}").Update(TestFirestoreI{}, testFsFunc)
		assert.NotNil(t, up.fn, "Firestore function should be not nil")

		err := up.reg.EntryPoint(testUpdatemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Write Exec", func(t *testing.T) {
		wr := reg.Firestore().Collection("testColl").Document("{uid}").Write(TestFirestoreI{}, testFsFunc)
		assert.NotNil(t, wr.fn, "Firestore function should be not nil")

		err := wr.reg.EntryPoint(testWritemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

}
