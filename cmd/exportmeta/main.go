// Command exportmeta writes the version-agnostic NCPDP model metadata
// (transaction types, segment/field codes, header layouts) as JSON so that
// non-Go consumers can drive their own serializers from it.
//
// Usage:
//
//	go run ./cmd/exportmeta [-out path]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/transactrx/NCPDPSerDe/pkg/metaexport"
)

func main() {
	out := flag.String("out", "", "output file path (default: stdout)")
	flag.Parse()

	meta, err := metaexport.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "exportmeta: %v\n", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "exportmeta: %v\n", err)
		os.Exit(1)
	}
	data = append(data, '\n')

	if *out == "" {
		os.Stdout.Write(data)
		return
	}

	if err := os.WriteFile(*out, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "exportmeta: %v\n", err)
		os.Exit(1)
	}
}
