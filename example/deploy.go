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
