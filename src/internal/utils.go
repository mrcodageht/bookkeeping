package internal

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"
)

var MustVerbose bool
var FileArg string
var DirArg string
var VendorsArg string
var FilterArg string
var ConfigArg string

func ReadFlag() {
	base := GetBase()
	flag.StringVar(&FileArg, "file", "", "Help message for file")
	flag.StringVar(&DirArg, "dir", "", "Help message for dir")
	flag.StringVar(&VendorsArg, "vendors", path.Join(base, "vendors.yaml"), "Fichier des vendors")
	flag.BoolVar(&MustVerbose, "verbose", false, "Activer la verbosite")
	flag.StringVar(&FilterArg, "filter", "", "Filter les transaction selon un vendor")
	flag.StringVar(&ConfigArg, "config", path.Join(base, "configs.yaml"), "Chemin du fichier configuration")
	flag.Parse()

	strings.TrimSpace(FileArg)
	strings.TrimSpace(DirArg)
	strings.TrimSpace(VendorsArg)
}

func GetBase() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}
