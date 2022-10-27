package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type FileOpts struct {
	Delimiter   rune
	ShowHeaders bool

	KeepCols []string
}

type File struct {
	Headers []string
	Lines   []Line
}

type Line []string

func HandleFile(opts FileOpts, file *os.File) (*File, error) {
	if file == nil {
		return nil, errors.New("missing file")
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	rdr.Comma = opts.Delimiter

	// Until we handle flagIndices or flagNotIndices we assume the first record contains headers
	headerIndexes := make(map[int]string)
	headers, err := rdr.Read()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	toKeep := opts.KeepCols
	for i := range headers {
		if len(toKeep) == 0 {
			headerIndexes[i] = strings.TrimSpace(headers[i])
		}

		for j := range toKeep {
			if strings.EqualFold(strings.TrimSpace(headers[i]), strings.TrimSpace(toKeep[j])) {
				headerIndexes[i] = strings.TrimSpace(headers[i])
			}
		}
	}

	var output File
	// TODO(adam): populate output.Headers

	lineNumber := 1
	for {
		lineNumber += 1

		cols, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("line %d failed to parse: %v", lineNumber, err)
		}

		var line Line
		for i := range cols {
			if _, exists := headerIndexes[i]; exists {
				line = append(line, cols[i])
			}
		}

		output.Lines = append(output.Lines, line)
	}

	return &output, nil
}
