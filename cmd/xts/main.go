package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/genkiroid/xts"
)

var (
	version string
	commit  string
	date    string
	builtBy string
)

func main() {
	var showVersion bool
	var yaml bool
	var outdir string

	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&yaml, "y", false, "output as yaml")
	flag.BoolVar(&yaml, "yaml", false, "output as yaml")
	flag.StringVar(&outdir, "d", "", "output to files at specified directory")
	flag.StringVar(&outdir, "dir", "", "output to files at specified directory")
	flag.Parse()

	if showVersion {
		fmt.Printf("xts, version %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "argument is missing")
		os.Exit(1)
	}

	path, err := filepath.Abs(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "filepath.Abs(%s) returned error: %s\n", flag.Args()[0], err)
		os.Exit(1)
	}

	input, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.ReadFile(%s) returned error: %s\n", path, err)
		os.Exit(1)
	}

	xml, err := xts.NewMySQLXMLDump(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewMySQLXMLDump(%s) returned error: %s\n", input, err)
		os.Exit(1)
	}

	sql := xts.NewSql(xml, outdir)
	if yaml {
		if err := sql.Yaml(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if err := sql.InsertStmt(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
