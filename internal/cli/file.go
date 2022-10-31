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
	Opts    FileOpts
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
	headerIndexes := make(map[int]int) // index within input file -> index in output
	headers, err := rdr.Read()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	output := File{
		Opts: opts,
	}

	toKeep := opts.KeepCols
	for headIdx := range headers {
		if len(toKeep) == 0 {
			headerIndexes[headIdx] = headIdx // keep all cols, so maintain order
			output.Headers = append(output.Headers, headers[headIdx])
		}

		for keepIdx := range toKeep {
			if strings.EqualFold(strings.TrimSpace(headers[headIdx]), strings.TrimSpace(toKeep[keepIdx])) {
				headerIndexes[headIdx] = keepIdx
				output.Headers = append(output.Headers, toKeep[keepIdx])
			}
		}
	}

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

		line := make(Line, len(toKeep))
		for i := range cols {
			if keepIdx, exists := headerIndexes[i]; exists {
				line[keepIdx] = cols[i] // format output according to -keep order
			}
		}

		output.Lines = append(output.Lines, line)
	}

	return &output, nil
}
