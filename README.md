# firebase-fx

## A wrapper for Google Cloud Functions that simplifies the deployment of serverless applications.

Meant to expose a similar API to the Firebase Functions package for node.js.

### Features
 - [ ] Deployment
    - [x] Output bash script
    - [ ] Profile memory for deployment: --memory flag
 - [ ] HTTP triggers
    - [ ] Unauthenticated
    - [ ] Methods, Headers, Host, Query
    - [ ] Middleware
 - [ ] Analytics
 - [x] Firebase Authentication triggers
 - [x] Firestore triggers
    - [x] Document path wildcards
      - [x] Access vars
    - [x] Custom data types
    - [ ] fx tagged fields
      - [x] string
      - [x] number = float64/int
      - [x] boolean = bool
      - [x] map = struct OR map[string]interface{}
        - [x] map = struct OR map[string]interface{}
      - [ ] array = []interface{}
      - [ ] geopoint = struct
      - [x] timestamp = time.Time
 - [x] PubSub triggers
    - [x] Custom data types
 - [x] Firebase Realtime Database triggers
    - [x] Path wildcards
      - [x] Access vars
    - [x] Custom data types - JSON tags
 - [ ] Schedule triggers
 - [ ] Storage triggers

 ### Usage


setup.go
```go
package main

import (
	"context"
	"fmt"

	register "github.com/schmorrison/firebase-fx"
)

var MyEntryPoint = register.SharedEntryPoint

func init() {
	register.Shared.PubSub("my-topic").Publish(MyCustomData{}, func(ctx context.Context, msg register.PubSubMessage) error {
		fmt.Println(msg.Topic)
		if data, ok := msg.Data.(*MyCustomData); ok {
			// do something with v
			fmt.Println(data)
		}
		return nil
	})

	register.Shared.Firestore().Collection("users").Document("{uid}").Create(MyUserData{}, func(ctx context.Context, e register.FirestoreEvent) error {
		fmt.Println(e.Vars()["uid"])

		if data, ok := e.Value.Fields.(*MyUserData); ok {
			// do something with v
			fmt.Println(data)
		}

		if data, ok := e.OldValue.Fields.(*MyUserData); ok {
			// do something with v
			fmt.Println(data)
		}
		return nil
	})
}

type MyCustomData struct {
	Name string
}

type MyUserData struct {
	Email string
}
```

deploy.go
```go
package main

import (
	"fmt"

	register "github.com/schmorrison/firebase-fx"
)

func main() {
	fmt.Println(register.Shared.
		WithEntrypoint("MyEntryPoint").
		WithProjectID("my-project-id").
		WithRuntime("go116").
		Verbosity(register.DebugVerbosity).
		Deploy(),
	)
}

```

command
```bash
go run setup.go deploy.go

gcloud functions deploy  --entry-point MyEntryPoint --runtime go116 --project my-project-id --verbosity debug \
pubsubpublish-my-topic --trigger-topic my-topic &&  \
gcloud functions deploy  --entry-point MyEntryPoint --runtime go116 --project my-project-id --verbosity debug \
firestoreDocCreate-users-uid --trigger-event providers/cloud.firestore/eventTypes/document.create --trigger-resource projects/my-project-id/databases/(default)/documents/users/{uid}
```
