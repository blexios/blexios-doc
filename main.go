package main

import (
	"os"

	doccmd "blexios/cmd/doc"
)

func main() {
	os.Exit(doccmd.RunGenerate(os.Args[1:]))
}
