//go:build js

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/adamdecaf/csvq/internal/cli"
	"github.com/adamdecaf/csvq/internal/format"
)

func main() {
	js.Global().Set("runner", runner())
	<-make(chan bool)
}

func runner() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return fmt.Sprintf("Got %d arguments, expected 1", len(args))
		}

		input := strings.NewReader(args[0].String())
		opts, err := parseFileOpts(args[1].String())
		if err != nil {
			return fmt.Sprintf("Unable to parse FileOpts: %v", err)
		}
		if opts == nil {
			return "Missing FileOpts"
		}

		file, err := cli.HandleFile(*opts, input)
		if err != nil {
			return fmt.Sprintf("Problem handling contents: %v", err)
		}

		var buf bytes.Buffer
		err = format.WriteFile(&buf, "csv", file)
		if err != nil {
			return fmt.Sprintf("Problem formatting output: %v", err)
		}
		return buf.String()
	})
}

func parseFileOpts(in string) (*cli.FileOpts, error) {
	var opts cli.FileOpts
	err := json.NewDecoder(strings.NewReader(in)).Decode(&opts)
	if err != nil {
		return nil, err
	}
	return &opts, nil
}
