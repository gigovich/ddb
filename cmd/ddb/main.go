package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gigovich/ddb/internal/parser"
)

func main() {
	flag.Parse()

	var file string
	if flag.NArg() > 0 {
		file = flag.Arg(0)
	}

	if file == "" {
		for _, e := range os.Environ() {
			if strings.HasPrefix(e, "GOFILE") {
				file = os.Getenv("GOFILE")
				break
			}
		}
	}

	if file == "" {
		fmt.Println("Usage: ddb [SOURCE_FILE]")
		fmt.Println("  [SOURCE_FILE] - file from which we parse model definition")
		fmt.Println("  NOTE: If [SOURCE_FILE] not set but $GOFILE environment value is defined,")
		fmt.Println("        then this value will be used as SOURCE_FILE,")
		fmt.Println("        useful for go:generate tag.")
	}

	ctx, err := parser.New(file).Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(ctx)
}
