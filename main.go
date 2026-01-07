package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type LineKind int

const (
	LineOther LineKind = iota
	LineAssign
)

type ParsedLine struct {
	Kind       LineKind
	Raw        string
	Prefix     string // "" or "export "
	Key        string
	Value      string
	Tail       string // includes spaces + comment (e.g. "  # foo")
	HasEquals  bool
	LineNumber int
}

type ParseResult struct {
	Lines       []ParsedLine
	ByKey       map[string]int
	KeysInOrder []string
}

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

	sourceBytes, err := os.ReadFile(sourceEnvPath)
	must(err)

	source := ParseDotenv(sourceBytes)
	fmt.Println(source)

	destinationBytes, err := os.ReadFile(destinationEnvPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			destinationBytes = []byte{}
		} else {
			must(err)
		}
	}

	destination := ParseDotenv(destinationBytes)
	fmt.Println(destination)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func RenderAssign(pl ParsedLine) string {
	return pl.Prefix + pl.Key + "=" + pl.Value + pl.Tail
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t'
}

// SplitValueAndTail splits RHS into value + tail (tail includes spaces + comment).
// It treats '#' as comment only when not inside quotes and preceded by whitespace.
func SplitValueAndTail(rhs string) (string, string) {
	inSingle := false
	inDouble := false
	escaped := false

	for i := 0; i < len(rhs); i++ {
		ch := rhs[i]

		if escaped {
			escaped = false
			continue
		}

		if inDouble && ch == '\\' {
			escaped = true
			continue
		}

		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if ch == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}

		if ch == '#' && !inSingle && !inDouble {
			if i == 0 || isSpace(rhs[i-1]) {
				start := i
				for start > 0 && isSpace(rhs[start-1]) {
					start--
				}
				value := rhs[:start]
				tail := rhs[start:]
				return strings.TrimRight(value, " \t"), tail
			}
		}
	}

	return strings.TrimRight(rhs, " \t"), ""
}

func ParseLine(raw string, lineNo int) ParsedLine {
	pl := ParsedLine{
		Kind:       LineOther,
		Raw:        raw,
		LineNumber: lineNo,
	}

	trimLeft := strings.TrimLeft(raw, " \t")
	if trimLeft == "" {
		return pl
	}
	if strings.HasPrefix(trimLeft, "#") {
		return pl
	}

	prefix := ""
	rest := trimLeft
	if strings.HasPrefix(rest, "export ") {
		prefix = "export "
		rest = strings.TrimLeft(rest[len("export "):], " \t")
	}

	eq := strings.IndexByte(rest, '=')
	if eq < 0 {
		return pl
	}

	key := strings.TrimSpace(rest[:eq])
	if key == "" {
		return pl
	}

	right := rest[eq+1:]
	value, tail := SplitValueAndTail(right)

	pl.Kind = LineAssign
	pl.Prefix = prefix
	pl.Key = key
	pl.Value = value
	pl.Tail = tail
	pl.HasEquals = true

	return pl
}

func ParseDotenv(b []byte) ParseResult {
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	res := ParseResult{
		ByKey: make(map[string]int),
	}

	lineNo := 0
	for sc.Scan() {
		lineNo++
		raw := sc.Text()

		pl := ParseLine(raw, lineNo)
		res.Lines = append(res.Lines, pl)

		if pl.Kind == LineAssign {
			// First occurrence wins. If you want "last wins", overwrite here.
			if _, exists := res.ByKey[pl.Key]; !exists {
				res.ByKey[pl.Key] = len(res.Lines) - 1
				res.KeysInOrder = append(res.KeysInOrder, pl.Key)
			}
		}
	}

	if err := sc.Err(); err != nil {
		must(err)
	}

	return res
}
