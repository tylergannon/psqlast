package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	pg_query "github.com/pganalyze/pg_query_go/v6"
	"github.com/tylerBrock/colorjson"
)

func main() {
	// Define flags
	noPretty := flag.Bool("no-pretty", false, "Disable pretty-printed (indented) JSON output")
	noColor := flag.Bool("no-color", false, "Disable colorized JSON output")
	outFile := flag.String("out", "", "Output file path")
	indentOpt := flag.String("indent", "tab", "Indentation style: 'tab' or an integer (e.g. '2', '4')")
	useClipboard := flag.Bool("use-clipboard", false, "Read SQL script from the system clipboard")

	// Alias for convenience: -c same as --use-clipboard
	flag.BoolVar(useClipboard, "c", false, "Read SQL script from the system clipboard (shorthand)")

	// Parse flags + leftover positional args
	flag.Parse()

	// Decide on input source
	var inputSQL []byte
	var err error

	if *useClipboard {
		// 1) If --use-clipboard/-c is set, read from clipboard
		cbContents, cbErr := clipboard.ReadAll()
		if cbErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to read from clipboard: %v\n", cbErr)
			os.Exit(1)
		}
		inputSQL = []byte(cbContents)

	} else {
		// 2) Otherwise, check if there's a positional file argument
		args := flag.Args()
		if len(args) > 1 {
			_, _ = fmt.Fprintln(os.Stderr, "Error: Too many positional arguments.")
			_, _ = fmt.Fprintln(os.Stderr, "Usage: [flags] [filePath]")
			os.Exit(1)
		} else if len(args) == 1 {
			filePath := args[0]
			// Validate file
			fi, errStat := os.Stat(filePath)
			if errStat != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Cannot access file '%s': %v\n", filePath, errStat)
				os.Exit(1)
			}
			if !fi.Mode().IsRegular() {
				_, _ = fmt.Fprintf(os.Stderr, "Error: '%s' is not a regular file.\n", filePath)
				os.Exit(1)
			}
			// Read file
			inputSQL, err = os.ReadFile(filePath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to read file '%s': %v\n", filePath, err)
				os.Exit(1)
			}

		} else {
			// 3) If no file argument, read from STDIN
			inputSQL, err = io.ReadAll(os.Stdin)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error: Failed to read from STDIN.")
				os.Exit(1)
			}
		}
	}

	// Convert SQL to JSON via pg_query
	data, err := pg_query.Parse(string(inputSQL))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to parse SQL: %v\n", err)
		os.Exit(1)
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to parse JSON: %v\n", err)
		os.Exit(1)
	}

	// Check if we need to pretty-print/indent:
	// "if the --no-pretty option is NOT given (or if --out is specified), the resulting JSON should be formatted / indented."
	needPretty := (!*noPretty) || (*outFile != "")

	if needPretty {
		// Try to unmarshal the JSON so we can re-format
		var obj interface{}
		if err = json.Unmarshal([]byte(jsonStr), &obj); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: Invalid JSON from parser: %v\n", err)
			os.Exit(1)
		}

		// Decide how many spaces per indent (or tab)
		var indentStr string
		spaces := 0
		if *indentOpt == "tab" {
			indentStr = "\t"
			spaces = -1 // marker for "tab"
		} else {
			n, convErr := strconv.Atoi(*indentOpt)
			if convErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Invalid value for --indent: '%s'\n", *indentOpt)
				os.Exit(1)
			}
			spaces = n
			indentStr = strings.Repeat(" ", n)
		}

		if !*noColor {
			// If color is enabled, we use the colorjson library
			// but note that colorjson only supports an integer indentation = number of spaces.
			// There's no native "tab" support, so we approximate a tab by using 1 space if user asked for tab.
			formatter := colorjson.NewFormatter()
			if spaces < 0 {
				// user asked for tabs => fallback to 1 space
				formatter.Indent = 1
			} else {
				formatter.Indent = spaces
			}

			coloredBytes, err := formatter.Marshal(obj)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to colorize JSON: %v\n", err)
				os.Exit(1)
			}
			jsonStr = coloredBytes
		} else {
			// No color => standard library indentation
			prettyBytes, err := json.MarshalIndent(obj, "", indentStr)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: JSON marshal indent failed: %v\n", err)
				os.Exit(1)
			}
			jsonStr = prettyBytes
		}

	} else {
		// If we do not need pretty-printing, we can still colorize if --no-color is not set
		if !*noColor {
			var obj interface{}
			if err = json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Invalid JSON from parser: %v\n", err)
				os.Exit(1)
			}
			formatter := colorjson.NewFormatter()
			// Minimal JSON => zero Indent
			formatter.Indent = 0

			coloredBytes, err := formatter.Marshal(obj)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to colorize JSON: %v\n", err)
				os.Exit(1)
			}
			jsonStr = coloredBytes
		}
	}

	// Now decide where to write output
	var writer io.Writer = os.Stdout
	if *outFile != "" {
		// Validate that we can open (create/truncate) the file
		f, createErr := os.Create(*outFile)
		if createErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to open output file '%s': %v\n", *outFile, createErr)
			os.Exit(1)
		}
		defer f.Close()
		writer = f
	}

	// Finally, print the JSON
	_, _ = fmt.Fprintln(writer, string(jsonStr))
}
