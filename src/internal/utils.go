package internal

import (
	"flag"
	"strings"
)

var MustVerbose bool
var FileArg string
var DirArg string
var VendorArg string
var FilterArg string

func ReadFlag() {
	flag.StringVar(&FileArg, "file", "", "Help message for file")
	flag.StringVar(&DirArg, "dir", "", "Help message for dir")
	flag.StringVar(&VendorArg, "vendors", "../bookkeping/vendors.yaml", "")
	flag.BoolVar(&MustVerbose, "verbose", false, "Activer la verbosite")
	flag.StringVar(&FilterArg, "filter", "", "Filter les transaction selon un vendor")
	flag.Parse()

	strings.TrimSpace(FileArg)
	strings.TrimSpace(DirArg)
	strings.TrimSpace(VendorArg)
}
