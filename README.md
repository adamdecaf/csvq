# csvq

[![GoDoc](https://godoc.org/github.com/adamdecaf/csvq?status.svg)](https://godoc.org/github.com/adamdecaf/csvq)
[![Build Status](https://github.com/adamdecaf/csvq/workflows/Go/badge.svg)](https://github.com/adamdecaf/csvq/actions)
[![Coverage Status](https://codecov.io/gh/adamdecaf/csvq/branch/master/graph/badge.svg)](https://codecov.io/gh/adamdecaf/csvq)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamdecaf/csvq)](https://goreportcard.com/report/github.com/adamdecaf/csvq)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/adamdecaf/csvq/master/LICENSE)

csvq is a CLI package for parsing and transforming CSV files. This is useful because often trimming down a CSV file can make processing it easier.

## Install

Download the [latest release for your architecture](https://github.com/adamdecaf/csvq/releases/latest).

You can install from source:
```
go install github.com/adamdecaf/csvq/cmd/csvq@latest
```

## Usage

Extract first_name and last_name columns (in that order). Sort results.
```
csvq -keep first_name,last_name ~/Downloads/report.csv | sort -u
```

Change delimiter used in `report.csv`.
```
csvq -d';' user_id,dob,email ~/Downloads/report.csv
```

Output CSV columns in a table.
```
csvq -keep first_name,last_name -format table
```

Combine multiple files.
```
csvq -keep user_id,email ~/Downloads/report1.csv ~/Downloads/report2.csv
```

## Supported and tested platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
