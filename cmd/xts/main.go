package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/genkiroid/xts"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Give the xml file path as an argument.")
	}

	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	input, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	xml, err := xts.NewMySQLXMLDump(input)
	if err != nil {
		log.Fatal(err)
	}

	sql := xts.NewSql(xml)
	fmt.Println(sql)
}
