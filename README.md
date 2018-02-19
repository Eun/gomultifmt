
# gomultifmt [![Travis](https://img.shields.io/travis/Eun/gomultifmt.svg)](https://travis-ci.org/Eun/gomultifmt) [![go-report](https://goreportcard.com/badge/github.com/Eun/gomultifmt)](https://goreportcard.com/report/github.com/Eun/gomultifmt)
Run multiple golang formatters in one command


## Installation
```
go get -u github.com/Eun/gomultifmt
```

## Usage
```
usage: gomultifmt [<flags>] [<path>...]

Run multiple golang formatters in one command

Flags:
  -h, --help             Show context-sensitive help (also try --help-long and --help-man).
  -f, --fmt=gofmt ...    Formatter to call (sepcify it multiple times, e.g.: --fmt=gofmt --fmt=goremovelines
  -w, --toSource         Write result to (source) file instead of stdout
  -s, --skip=DIR... ...  Skip directories with this name when expanding '...'.
      --vendor           Enable vendoring support (skips 'vendor' directories and sets GO15VENDOREXPERIMENT=1).
  -d, --debug            Display debug messages.
  -v, --version          Show application version.

Args:
  [<path>]  Directories to format. Defaults to ".". <path>/... will recurse.
```


## Examples
### Format all files with gofmt first, and then goremovelines afterwards
```
gomultifmt --vendor --fmt=gofmt --fmt=goremovelines -w ./...
```

### Format main.go with gofmt and goremovelines, but write to stdout

```
gomultifmt --fmt=gofmt --fmt=goremovelines main.go
```

## Bonus VSCode config

  "go.formatFlags": [
      "--fmt=goreturns",
      "--fmt=goremovelines",
  ],
  "go.formatTool": "gomultifmt",

