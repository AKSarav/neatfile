package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	version    = flag.Bool("v", false, "Print version")
	help       = flag.Bool("h", false, "Show help")
	output     = flag.String("o", "", "Output file (default: stdout)")
	keepEmpty  = flag.Bool("k", false, "Keep empty lines in output")
	versiontag = "0.0.1"
)

// Map of file extensions to comment patterns
var fileCommentPatterns = map[string][]*regexp.Regexp{
	".yaml":   {regexp.MustCompile(`^\s*#.*`)},                                        // YAML
	".yml":    {regexp.MustCompile(`^\s*#.*`)},                                        // YAML
	".json":   {regexp.MustCompile(`^\s*//.*`)},                                       // JSON
	".ini":    {regexp.MustCompile(`^\s*;.*`)},                                        // INI
	".sql":    {regexp.MustCompile(`^\s*--.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // SQL
	".js":     {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // JavaScript
	".ts":     {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // TypeScript
	".jsx":    {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // JSX
	".tsx":    {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // TypeScript
	".go":     {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // Go
	".py":     {regexp.MustCompile(`^\s*#.*`), regexp.MustCompile(`(?s)""".*?"""`)},   // Python
	".rb":     {regexp.MustCompile(`^\s*#.*`), regexp.MustCompile(`(?s)""".*?"""`)},   // Ruby
	".sh":     {regexp.MustCompile(`^\s*#.*`)},                                        // Shell
	".pl":     {regexp.MustCompile(`^\s*#.*`)},                                        // Perl
	".php":    {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // PHP
	".java":   {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // Java
	".c":      {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // C
	".h":      {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // C
	".cpp":    {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // C++
	".hpp":    {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // C++
	".cs":     {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // C#
	".rs":     {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // Rust
	".swift":  {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // Swift
	".kt":     {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)/\*.*?\*/`)},  // Kotlin
	".clj":    {regexp.MustCompile(`^\s*;.*`)},                                        // Clojure
	".cljs":   {regexp.MustCompile(`^\s*;.*`)},                                        // ClojureScript
	".cljc":   {regexp.MustCompile(`^\s*;.*`)},                                        // Clojure
	".edn":    {regexp.MustCompile(`^\s*;.*`)},                                        // Clojure
	".lisp":   {regexp.MustCompile(`^\s*;.*`)},                                        // Common Lisp
	".rkt":    {regexp.MustCompile(`^\s*;.*`)},                                        // Racket
	".scm":    {regexp.MustCompile(`^\s*;.*`)},                                        // Scheme
	".ss":     {regexp.MustCompile(`^\s*;.*`)},                                        // Scheme
	".el":     {regexp.MustCompile(`^\s*;.*`)},                                        // Emacs Lisp
	".ex":     {regexp.MustCompile(`^\s*#.*`)},                                        // Elixir
	".exs":    {regexp.MustCompile(`^\s*#.*`)},                                        // Elixir
	".erl":    {regexp.MustCompile(`^\s*%.*`)},                                        // Erlang
	".hrl":    {regexp.MustCompile(`^\s*%.*`)},                                        // Erlang
	".hs":     {regexp.MustCompile(`^\s*--.*`), regexp.MustCompile(`(?s)\{-.*-\}`)},   // Haskell
	".lhs":    {regexp.MustCompile(`^\s*--.*`), regexp.MustCompile(`(?s)\{-.*-\}`)},   // Haskell
	".ml":     {regexp.MustCompile(`^\s*\(\*.*\*\)`), regexp.MustCompile(`^\s*//.*`)}, // OCaml
	".mli":    {regexp.MustCompile(`^\s*\(\*.*\*\)`), regexp.MustCompile(`^\s*//.*`)}, // OCaml
	".fs":     {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)\(\*.*\*\)`)}, // F#
	".fsi":    {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)\(\*.*\*\)`)}, // F#
	".fsx":    {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)\(\*.*\*\)`)}, // F#
	".fsproj": {regexp.MustCompile(`^\s*//.*`), regexp.MustCompile(`(?s)\(\*.*\*\)`)}, // F#
}

func getCommentPatterns(filename string) []*regexp.Regexp {
	ext := filepath.Ext(filename)
	if patterns, found := fileCommentPatterns[ext]; found {
		return patterns
	}
	return []*regexp.Regexp{} // Default to no removal if unknown file type
}

func isCommentOrEmpty(line string, patterns []*regexp.Regexp) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" && !*keepEmpty {
		return true
	}
	for _, pattern := range patterns {
		if pattern.MatchString(trimmed) {
			return true
		}
	}
	return false
}

func processFile(filename string, writer *bufio.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	patterns := getCommentPatterns(filename)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !isCommentOrEmpty(line, patterns) {
			writer.WriteString(line + "\n")
		}
	}
	return scanner.Err()
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: neatfile [-v] [-k] [-o output] <file>
Options:
  -h    Show help
  -v    Print version
  -k    Keep empty lines in output (default is to remove them)
  -o    Output file (default: stdout)
`)
}
func main() {
	flag.Parse()
	file := flag.Args()

	if *help {
		usage()
		os.Exit(0)
	}

	if *version {
		fmt.Println("neatfile ", versiontag)
		os.Exit(0)
	}

	if len(file) == 0 && !*version {
		usage()
		os.Exit(1)
	} else if len(file) > 1 {
		fmt.Fprintf(os.Stderr, "Error: only one file is allowed at a time \n")
		os.Exit(1)
	}

	for _, f := range file {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: file %s does not exist\n", f)
			os.Exit(1)
		} else if info, err := os.Stat(f); err == nil && info.IsDir() {
			fmt.Fprintf(os.Stderr, "Error: %s is a directory\n", f)
			os.Exit(1)
		} else if _, err := os.Open(f); err != nil {
			fmt.Fprintf(os.Stderr, "Error: file is present but we could not open - %s \ncould it be permission issue ?\n", f)
			os.Exit(1)
		}
	}

	var writer *bufio.Writer
	if *output != "" {
		outFile, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer outFile.Close()
		writer = bufio.NewWriter(outFile)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}

	for _, f := range file {
		err := processFile(f, writer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", file, err)
		}
	}

	writer.Flush()
}
