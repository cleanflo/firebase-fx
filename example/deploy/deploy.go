package main

import (
	"fmt"

	register "github.com/cleanflo/firebase-fx"
	functions "github.com/cleanflo/firebase-fx/functions"
)

func main() {
	fmt.Println(functions.Register.
		WithEntrypoint("Register.EntryPoint").
		WithProjectID("my-project-id").
		WithRuntime("go116").
		Verbosity(register.DebugVerbosity).
		Deploy(),
	)
}
