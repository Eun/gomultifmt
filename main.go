package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	fmtFlag           = kingpin.CommandLine.Flag("fmt", "Formatter to call (sepcify it multiple times, e.g.: -f=gofmt -f=goremovelines").Short('f').PlaceHolder("gofmt").Default("gofmt").Strings()
	writeToSourceFlag = kingpin.CommandLine.Flag("toSource", "Write result to (source) file instead of stdout").Short('w').Default("false").Bool()
	skipFlag          = kingpin.CommandLine.Flag("skip", "Skip directories with this name when expanding '...'.").Short('s').PlaceHolder("DIR...").Strings()
	vendorFlag        = kingpin.CommandLine.Flag("vendor", "Enable vendoring support (skips 'vendor' directories and sets GO15VENDOREXPERIMENT=1).").Bool()
	debugFlag         = kingpin.CommandLine.Flag("debug", "Display debug messages.").Short('d').Bool()
)

type tool struct {
	cmd  string
	args []string
}

func formatPaths(paths []string, tools []tool) {
	for i := 0; i < len(paths); i++ {
		for j := 0; j < len(tools); j++ {
			out := &bytes.Buffer{}
			args := append(tools[j].args, paths[i])
			debug("Running tool `%s' with `%v'", tools[j].cmd, args)

			cmd := exec.Command(tools[j].cmd, args...)

			cmd.Stdout = out
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				warning("Tool `%s' failed: %v", tools[j].cmd, err)
				continue
			}

			if writeToSourceFlag != nil && *writeToSourceFlag {
				f, err := os.Create(paths[i])
				if err == nil {
					if _, err = f.Write(out.Bytes()); err != nil {
						warning("Unable to write file `%s': %v", paths[i], err)
					}
					if err = f.Close(); err != nil {
						warning("Unable to close file `%s': %v", paths[i], err)
					}
				} else {
					warning("Unable to create file `%s': %v", paths[i], err)
				}
			} else {
				if _, err := io.Copy(os.Stdout, out); err != nil {
					warning("Unable to write to stdout (`%s'): %v", paths[i], err)
				}
			}
		}
	}
}

func main() {
	pathsArg := kingpin.Arg("path", "Directories to format. Defaults to \".\". <path>/... will recurse.").Strings()
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.Version("gomultifmt 1.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.Help = "Run multiple golang formatters in one command"

	kingpin.Parse()

	if pathsArg == nil || len(*pathsArg) <= 0 {
		return
	}

	if fmtFlag == nil || len(*fmtFlag) <= 0 {
		return
	}

	if skipFlag == nil {
		skipFlag = &[]string{}
	}

	if os.Getenv("GO15VENDOREXPERIMENT") == "1" || (vendorFlag != nil && *vendorFlag) {
		if err := os.Setenv("GO15VENDOREXPERIMENT", "1"); err != nil {
			warning("setenv GO15VENDOREXPERIMENT: %s", err)
		}
		*skipFlag = append(*skipFlag, "vendor")
		trueValue := true
		vendorFlag = &trueValue
	}

	var tools []tool

	for i := 0; i < len(*fmtFlag); i++ {
		var tool tool
		parts := strings.Split((*fmtFlag)[i], " ")
		for j := 0; j < len(parts); j++ {
			part := strings.TrimSpace(parts[j])
			if len(part) > 0 {
				if len(tool.cmd) <= 0 {
					tool.cmd = part
				} else {
					tool.args = append(tool.args, part)
				}
			}
		}
		if len(tool.cmd) > 0 {
			tools = append(tools, tool)
		}
	}

	formatPaths(resolvePaths(*pathsArg, *skipFlag), tools)
}
