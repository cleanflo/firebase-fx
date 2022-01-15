package register

import (
	"context"
	"encoding/json"
	"testing"

	"cloud.google.com/go/functions/metadata"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	reg := NewRegister()

	testArchivemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(StorageObjectArchiveEvent),
		Resource: &metadata.Resource{
			Name: "projects/_/buckets/testBucket/objects/profile/image.jpg",
		},
	})
	testDeletemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(StorageObjectDeleteEvent),
		Resource: &metadata.Resource{
			Name: "projects/_/buckets/testBucket/objects/profile/image.jpg",
		},
	})
	testFinalizemd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(StorageObjectFinalizeEvent),
		Resource: &metadata.Resource{
			Name: "projects/_/buckets/testBucket/objects/profile/image.jpg",
		},
	})
	testMetamd := metadata.NewContext(context.Background(), &metadata.Metadata{
		EventType: string(StorageObjectMetadataUpdateEvent),
		Resource: &metadata.Resource{
			Name: "projects/_/buckets/testBucket/objects/profile/image.jpg",
		},
	})

	testDec := &Decoder{}
	err := json.Unmarshal([]byte(`{
		"kind":"storage#object",
		"id":"cleanflo-test-bucket/1 S Morrison.pdf/1642210177215991",
		"selfLink":"https://www.googleapis.com/storage/v1/b/cleanflo-test-bucket/o/1%20S%20Morrison.pdf",
		"name":"1 S Morrison.pdf",
		"bucket":"cleanflo-test-bucket",
		"generation":"1642210177215991",
		"metageneration":"3",
		"contentType":"application/pdf",
		"timeCreated":"2022-01-15T01:29:37.279Z",
		"updated":"2022-01-15T01:39:30.599Z",
		"temporaryHold":false,
		"eventBasedHold":false,
		"storageClass":"STANDARD",
		"timeStorageClassUpdated":"2022-01-15T01:29:37.279Z",
		"size":"79734",
		"md5Hash":"+ndvErmbIgA5ccSgYOCzxg==",
		"mediaLink":"https://www.googleapis.com/download/storage/v1/b/cleanflo-test-bucket/o/1%20S%20Morrison.pdf?generation=1642210177215991&alt=media",
		"metadata":{"ttt":"123"},
		"crc32c":"rZYNFQ==",
		"etag":"CPfbidLNsvUCEAM=",
		"customTime":"2022-01-18T06:00:00.000Z"
	}`), testDec)
	if err != nil {
		t.Errorf("Error unmarshalling test storage data: %v", err)
	}

	testStorageFunc := func(ctx context.Context, e StorageEvent) error {
		return nil
	}

	t.Run("Bucket", func(t *testing.T) {
		st := reg.Storage()
		st.Bucket("testRef")
		assert.Equalf(t, "testRef", st.resource, "path should be testRef, got: %s", st.resource)

		st.Bucket("sub")
		assert.Equalf(t, "sub", st.resource, "path should be sub, got: %s", st.resource)

		st.Bucket("deep/sub")
		assert.Equalf(t, "deep/sub", st.resource, "path should be deep/sub, got: %s", st.resource)
	})

	t.Run("Register Archive", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").Archive(testStorageFunc)
		assert.Same(t, st, reg.findStorage(StorageObjectArchiveEvent, "projects/_/buckets/testBucket/objects/profile/image.jpg"), "Storage function should be registered")
		assert.Same(t, st, reg.storage[StorageObjectArchiveEvent]["testBucket"], "Storage function should be registered")
		assert.NotNil(t, st.fn, "Storage function should be equal not nil")

		t.Log("Storage Function registered for ArchiveEvent")
	})

	t.Run("Register Delete", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").Delete(testStorageFunc)
		assert.Same(t, st, reg.findStorage(StorageObjectDeleteEvent, "projects/_/buckets/testBucket/objects/profile/image.jpg"), "Storage function should be registered")
		assert.Same(t, st, reg.storage[StorageObjectDeleteEvent]["testBucket"], "Storage function should be registered")
		assert.NotNil(t, st.fn, "Storage function should be equal not nil")

		t.Log("Storage Function registered for ArchiveEvent")
	})

	t.Run("Register Finalize", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").Finalize(testStorageFunc)
		assert.Same(t, st, reg.findStorage(StorageObjectFinalizeEvent, "projects/_/buckets/testBucket/objects/profile/image.jpg"), "Storage function should be registered")
		assert.Same(t, st, reg.storage[StorageObjectFinalizeEvent]["testBucket"], "Storage function should be registered")
		assert.NotNil(t, st.fn, "Storage function should be equal not nil")

		t.Log("Storage Function registered for ArchiveEvent")
	})

	t.Run("Register MetadataUpdate", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").MetadataUpdate(testStorageFunc)
		assert.Same(t, st, reg.findStorage(StorageObjectMetadataUpdateEvent, "projects/_/buckets/testBucket/objects/profile/image.jpg"), "Storage function should be registered")
		assert.Same(t, st, reg.storage[StorageObjectMetadataUpdateEvent]["testBucket"], "Storage function should be registered")
		assert.NotNil(t, st.fn, "Storage function should be equal not nil")

		t.Log("Storage Function registered for ArchiveEvent")
	})

	t.Run("Archive Exec", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").Archive(testStorageFunc)
		assert.NotNil(t, st.fn, "Storage function should be not nil")

		err := st.reg.EntryPoint(testArchivemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Delete Exec", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").Delete(testStorageFunc)
		assert.NotNil(t, st.fn, "Storage function should be not nil")

		err := st.reg.EntryPoint(testDeletemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("Finalize Exec", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").Finalize(testStorageFunc)
		assert.NotNil(t, st.fn, "Storage function should be not nil")

		err := st.reg.EntryPoint(testFinalizemd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

	t.Run("MetadataUpdate Exec", func(t *testing.T) {
		st := reg.Storage().Bucket("testBucket").MetadataUpdate(testStorageFunc)
		assert.NotNil(t, st.fn, "Storage function should be not nil")

		err := st.reg.EntryPoint(testMetamd, testDec)
		assert.Nil(t, err, "Error should be nil")
	})

}
