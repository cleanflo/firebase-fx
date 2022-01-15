package register

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploy(t *testing.T) {
	reg := NewRegister()

	t.Run("Auth Register", func(t *testing.T) {
		create := reg.Authentication().Create(nil)
		assert.Nil(t, create.fn, "Authentication function should be nil")
		delete := reg.Authentication().Delete(nil)
		assert.Nil(t, delete.fn, "Authentication function should be nil")
	})

	t.Run("Firestore Register", func(t *testing.T) {
		cr := reg.Firestore().Collection("testColl").Document("{uid}").Create(TestFirestoreI{}, nil)
		assert.Nil(t, cr.fn, "Firestore function should be nil")
		de := reg.Firestore().Collection("testColl").Document("{uid}").Delete(TestFirestoreI{}, nil)
		assert.Nil(t, de.fn, "Firestore function should be nil")
		up := reg.Firestore().Collection("testColl").Document("{uid}").Update(TestFirestoreI{}, nil)
		assert.Nil(t, up.fn, "Firestore function should be nil")
		wr := reg.Firestore().Collection("testColl").Document("{uid}").Write(TestFirestoreI{}, nil)
		assert.Nil(t, wr.fn, "Firestore function should be nil")
	})

	t.Run("PubSub Register", func(t *testing.T) {
		ps := reg.PubSub("test-topic").Publish(TestPubSubI{}, nil)
		assert.Nil(t, ps.fn, "PubSub function should be nil")
	})

	t.Run("RTDB Register", func(t *testing.T) {
		cr := reg.RealtimeDB().Ref("testColl/{uid}").Create(TestRTDBI{}, nil)
		assert.Nil(t, cr.fn, "RealtimeDB function should be not nil")

		de := reg.RealtimeDB().Ref("testColl/{uid}").Delete(TestRTDBI{}, nil)
		assert.Nil(t, de.fn, "RealtimeDB function should be not nil")

		up := reg.RealtimeDB().Ref("testColl/{uid}").Update(TestRTDBI{}, nil)
		assert.Nil(t, up.fn, "RealtimeDB function should be not nil")

		wr := reg.RealtimeDB().Ref("testColl/{uid}").Write(TestRTDBI{}, nil)
		assert.Nil(t, wr.fn, "RealtimeDB function should be not nil")
	})

	t.Run("Storage Register", func(t *testing.T) {
		ar := reg.Storage().Bucket("testBucket").Archive(nil)
		assert.Nil(t, ar.fn, "Storage function should be nil")

		de := reg.Storage().Bucket("testBucket").Delete(nil)
		assert.Nil(t, de.fn, "Storage function should be nil")

		fl := reg.Storage().Bucket("testBucket").Finalize(nil)
		assert.Nil(t, fl.fn, "Storage function should be nil")

		mu := reg.Storage().Bucket("testBucket").MetadataUpdate(nil)
		assert.Nil(t, mu.fn, "Storage function should be nil")
	})

	t.Run("Deploy", func(t *testing.T) {
		fmt.Println(reg.Deploy())
	})
}
