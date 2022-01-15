package register

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/functions/metadata"
)

const gcsBasePath = "projects/_/buckets/*/objects"

// Storage returns a new StorageFunction with the FunctionRegistrar set to the parent
func (f *FunctionRegistrar) Storage() *StorageFunction {
	s := &StorageFunction{reg: f}
	return s
}

func (f *FunctionRegistrar) findStorage(event StorageEventType, ref string) *StorageFunction {
	if f.storage[event] != nil && len(f.storage[event]) > 0 {
		refParts := strings.Split(ref, "/")
		if len(refParts) > 5 {
			bucket := refParts[3]
			if s := f.storage[event][bucket]; s != nil {
				return f.storage[event][bucket]
			}
		}
	}

	return nil
}

// StorageFunction  is a wrapper for the expected data, StorageFunc and the parentFunctionRegistrar
// Implements the CloudEventFunction interface
type StorageFunction struct {
	cloudDeployer
	reg *FunctionRegistrar
	fn  StorageFunc
}

// StorageFunc is the function signature for Google Cloud Storage Cloud Events
type StorageFunc func(ctx context.Context, e StorageEvent) error

// Bucket sets the given path to the bucket that the StorageFunc is executed on
func (s *StorageFunction) Bucket(path string) *StorageFunction {
	s.resource = path
	return s
}

// Finalize registers the specified function to the ObjectFinalizeEvent for Storage CloudEvent
//google.storage.object.finalize
func (s *StorageFunction) Finalize(fn StorageFunc) *StorageFunction {
	s.fn = fn

	if s.reg.storage[StorageObjectFinalizeEvent] == nil {
		s.reg.storage[StorageObjectFinalizeEvent] = make(map[string]*StorageFunction)
	}

	s.reg.storage[StorageObjectFinalizeEvent][s.resource] = s

	s.event = StorageObjectFinalizeEvent
	s.reg.events[s.Name()] = s
	return s
}

// Delete registers the specified function to the ObjectDeleteEvent for Storage CloudEvent
//google.storage.object.delete
func (s *StorageFunction) Delete(fn StorageFunc) *StorageFunction {
	s.fn = fn

	if s.reg.storage[StorageObjectDeleteEvent] == nil {
		s.reg.storage[StorageObjectDeleteEvent] = make(map[string]*StorageFunction)
	}

	s.reg.storage[StorageObjectDeleteEvent][s.resource] = s

	s.event = StorageObjectDeleteEvent
	s.reg.events[s.Name()] = s
	return s
}

// Archive registers the specified function to the ObjectArchiveEvent for Storage CloudEvent
//google.storage.object.archive
func (s *StorageFunction) Archive(fn StorageFunc) *StorageFunction {
	s.fn = fn

	if s.reg.storage[StorageObjectArchiveEvent] == nil {
		s.reg.storage[StorageObjectArchiveEvent] = make(map[string]*StorageFunction)
	}

	s.reg.storage[StorageObjectArchiveEvent][s.resource] = s

	s.event = StorageObjectArchiveEvent
	s.reg.events[s.Name()] = s
	return s
}

// MetadataUpdate registers the specified function to the ObjectMetadataUpdateEvent for Storage CloudEvent
//google.storage.object.metadataUpdate
func (s *StorageFunction) MetadataUpdate(fn StorageFunc) *StorageFunction {
	s.fn = fn

	if s.reg.storage[StorageObjectMetadataUpdateEvent] == nil {
		s.reg.storage[StorageObjectMetadataUpdateEvent] = make(map[string]*StorageFunction)
	}

	s.reg.storage[StorageObjectMetadataUpdateEvent][s.resource] = s

	s.event = StorageObjectMetadataUpdateEvent
	s.reg.events[s.Name()] = s
	return s
}

// GCSEvent is the expected payload for Google Cloud Storage CloudEvents.
/* Finalize:
{
	"kind":"storage#object",
	"id":"cleanflo-test-bucket/1 S Morrison.pdf/1642210177215991",
	"selfLink":"https://www.googleapis.com/storage/v1/b/cleanflo-test-bucket/o/1%20S%20Morrison.pdf",
	"name":"1 S Morrison.pdf",
	"bucket":"cleanflo-test-bucket",
	"generation":"1642210177215991",
	"metageneration":"1",
	"contentType":"application/pdf",
	"timeCreated":"2022-01-15T01:29:37.279Z",
	"updated":"2022-01-15T01:29:37.279Z"
	"storageClass":"STANDARD",
	"timeStorageClassUpdated":"2022-01-15T01:29:37.279Z",
	"size":"79734",
	"md5Hash":"+ndvErmbIgA5ccSgYOCzxg==",
	"mediaLink":"https://www.googleapis.com/download/storage/v1/b/cleanflo-test-bucket/o/1%20S%20Morrison.pdf?generation=1642210177215991&alt=media",
	"crc32c":"rZYNFQ==",
	"etag":"CPfbidLNsvUCEAE=",
}

MetadataUpdate:
{
	"kind":"storage#object",
	"id":"cleanflo-test-bucket/1 S Morrison.pdf/1642210177215991",
	"selfLink":"https://www.googleapis.com/storage/v1/b/cleanflo-test-bucket/o/1%20S%20Morrison.pdf",
	"name":"1 S Morrison.pdf",
	"bucket":"cleanflo-test-bucket",
	"generation":"1642210177215991",
	"metageneration":"3",
	"contentType":"application/pdf",
	"timeCreated":"2022-01-15T01:29:37.279Z",
	"updated":"2022-01-15T01:39:30.599Z"
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
	"customTime":"2022-01-18T06:00:00.000Z",
}

Delete:

{
	"kind":"storage#object",
	"id":"cleanflo-test-bucket/1 S Morrison.pdf/1642210177215991",
	"selfLink":"https://www.googleapis.com/storage/v1/b/cleanflo-test-bucket/o/1%20S%20Morrison.pdf",
	"name":"1 S Morrison.pdf",
	"bucket":"cleanflo-test-bucket",
	"generation":"1642210177215991",
	"metageneration":"3",
	"contentType":"application/pdf",
	"timeCreated":"2022-01-15T01:29:37.279Z",
	"updated":"2022-01-15T01:39:30.599Z"
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
	"customTime":"2022-01-18T06:00:00.000Z",
}
*/
type StorageEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

// CloudEventFunction

// HandleCloudEvent handles the Google Cloud Storage CloudEvent and calls the registered AuthenticationFunc
func (a *StorageFunction) HandleCloudEvent(ctx context.Context, md *metadata.Metadata, dec *Decoder) error {
	event := StorageEvent{}

	err := dec.Decode(&event)
	if err != nil {
		return Debug.Errf("failed to decode realtimeDB event [%s]: %s: %s", md.EventType, err, string(dec.data))
	}

	err = a.fn(ctx, event)
	if err != nil {
		return Debug.Errf("registered realtimeDBFunc failed [%s]: %s: RealtimeDBFunc %+v", md.EventType, err, a)
	}

	return nil
}

// Name returns the name of the function: "storageObject{Archive,Delete,Finalize,Metadata}-{bucketref}"
func (a *StorageFunction) Name() string {
	return strings.ToLower(fmt.Sprintf("%s-%s", a.event, a.Resource()))
}

// Resource returns the resource of the function: "{bucketRef}"
func (a *StorageFunction) Resource() string {
	return a.resource
}

// Event returns the EventType of the function:
//  StorageObjectFinalizeEvent / StorageObjectDeleteEvent / StorageObjectArchiveEvent / StorageObjectMetadataUpdateEvent / EventTypeUserDelete
func (a *StorageFunction) Event() EventType {
	return a.event.Type()
}
