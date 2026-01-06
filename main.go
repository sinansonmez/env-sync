package main

import (
	"flag"
	"fmt"
)

func main() {
	var sourceEnvPath string
	var destinationEnvPath string
	var dryRun bool
	var keepUnused bool
	var useSourceDefaults bool
	var fillEmpty bool

	flag.StringVar(&sourceEnvPath, "source", ".env.uat", "path to source env file (uat/test/dev)")
	flag.StringVar(&destinationEnvPath, "dest", ".env.prod", "path to destination env file (prod)")
	flag.BoolVar(&dryRun, "dry-run", false, "print result; do not write destination")
	flag.BoolVar(&keepUnused, "keep-unused", true, "append keys found only in destination to the end of output")
	flag.BoolVar(&useSourceDefaults, "use-source-defaults", false, "when a key is missing in destination, keep the default value from source instead of blank")
	flag.BoolVar(&fillEmpty, "fill-empty", false, "if a key exists in destination but value is empty, fill from source (still does not overwrite non-empty)")
	flag.Parse()

	fmt.Println("Hello, World!")
}
