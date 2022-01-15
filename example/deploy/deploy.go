package main

import (
	"fmt"

	functions "github.com/cleanflo/firebase-fx/functions"
	register "github.com/schmorrison/firebase-fx"
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
