package main

import (
	"flag"
	"fmt"
	"importtidier/logic"
	"os"
	"strings"
)

// This script is to tidy import groups of a project. It is a complementary tool
// for existing goimports and gofmt. Currently goimports would also tidy imports
// and arrange the import but it is not performing as expected in below cases:
//
// 1. If there is line break between standard libraries, it will leave the line break there
// 2. It will treat local packages(such as private repo in an org) as standard libs
//
// Usage:
//
//	go run main.go -f <filename> [-l <local_pkg1,local_pkg2>]
//
// Examples:
//
//	go run main.go -f /user/some/go/src/util.go -l private.example.com,localpg
//
// This script can be easily configured in GoLand as a File Watcher so that the import
// groups can be tidied automatically while saving the file
//
// This script is recreated based on https://github.com/corsc/go-tools/blob/master/fiximports
func main() {
	filename, localPkgStr := "", ""

	flag.Usage = usage
	flag.StringVar(&filename, "f", "", "the file to be tidied")
	flag.StringVar(&localPkgStr, "l", "", "local pkgs, delimited by comma, e.g: local.com,localB")
	flag.Parse()

	if filename == "" {
		usage()
		os.Exit(0)
	}

	var localPkgs []string
	if strings.TrimSpace(localPkgStr) != "" {
		localPkgs = strings.Split(strings.TrimSpace(localPkgStr), ",")
	}
	logic.ProcessFile(filename, localPkgs)
}

func usage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Println("go run main.go -f <filename> [-l <local_pkgs>]")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}
