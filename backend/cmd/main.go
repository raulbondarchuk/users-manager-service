package main

import (
	"app/cmd/build"
	"app/internal/composition"
	"flag"
)

func main() {
	// Define una bandera para el comando de compilación
	buildFlag := flag.Bool("build", false, "Ejecutar el proceso de compilación")
	flag.Parse()

	if *buildFlag {
		build.Build()
	} else {
		composition.Run()
	}
}
