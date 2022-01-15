package functions

import (
	"context"
	"fmt"

	register "github.com/schmorrison/firebase-fx"
)

var Register = register.Shared

func init() {
	Register.PubSub("my-topic").Publish(MyCustomData{}, func(ctx context.Context, msg register.PubSubMessage) error {
		fmt.Println(msg.Topic)
		if data, ok := msg.Data.(*MyCustomData); ok {
			// do something with v
			fmt.Println(data)
		}
		return nil
	})

	Register.Firestore().Collection("users").Document("{uid}").Create(MyUserData{}, func(ctx context.Context, e register.FirestoreEvent) error {
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
