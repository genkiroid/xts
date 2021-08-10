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

	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("xts, version %s\n", version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: xts path-to-mysql-xml-dump-file")
		os.Exit(1)
	}

	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "filepath.Abs(%s) returned error: %s\n", os.Args[1], err)
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

	sql := xts.NewSql(xml)
	fmt.Println(sql)
	os.Exit(0)
}
