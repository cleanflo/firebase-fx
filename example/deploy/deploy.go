package main

import (
	"fmt"

	register "github.com/cleanflo/firebase-fx"
	functions "github.com/cleanflo/firebase-fx/example"
)

func main() {
	fmt.Println(functions.Register.
		WithRegistrar("Register").
		WithProjectID("my-project-id").
		WithRuntime("go116").
		Verbosity(register.DebugVerbosity).
		Deploy(),
	)
}
