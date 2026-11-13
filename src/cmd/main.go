package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"mrccommunity.com/bookkeping/src/internal"
	"mrccommunity.com/bookkeping/src/internal/matcher"
	"mrccommunity.com/bookkeping/src/internal/model"
	"mrccommunity.com/bookkeping/src/internal/parser"
)

func main() {

	base := internal.GetBase()

	fmt.Println(base)
	internal.ReadFlag()
	if len(internal.FileArg) == 0 && len(internal.DirArg) == 0 {
		synopsis()
	}

	isDir := len(internal.DirArg) > 0

	if internal.MustVerbose {
		if !isDir {
			println("file :", internal.FileArg)
		} else {
			println("dir :", internal.DirArg)
		}
		if len(internal.FilterArg) > 0 {
			println("filter :", internal.FilterArg)
		}
	}
	var transactions []*model.Transaction
	var err error
	var vendors *model.VendorsYAML
	var configs *model.Config

	configs, err = internal.LoadConfig(internal.ConfigArg)
	logAndFail(err)

	vendors, err = internal.LoadVendor(internal.VendorsArg)
	logAndFail(err)

	if internal.MustVerbose {
		fmt.Printf("--- configs:\n%v\n\n", configs)
		log.Println("--- vendors")
		for _, v := range vendors.Vendors {
			fmt.Println((&v).ToString())
		}
	}

	if !isDir {
		transactions, err = parser.Parser(internal.FileArg, configs.Output.ErrorFile)
	} else {
		entry, err := os.ReadDir(internal.DirArg)
		logAndFail(err)
		tTemp := make([]*model.Transaction, 1)
		for _, e := range entry {
			path := path.Join(internal.DirArg, e.Name())
			if internal.MustVerbose {
				log.Printf("Found %v\n", path)
			}
			tTemp, err = parser.Parser(path, configs.Output.ErrorFile)

			transactions = append(transactions, tTemp...)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	matcher.Match(transactions, vendors)

	if !isBlank(internal.FilterArg) {
		if internal.MustVerbose {
			log.Printf("Applying filter [%s]", internal.FilterArg)
		}
		transactions = filter(transactions, internal.FilterArg)
	}
	nbMatched := 0
	nbUnknown := 0
	totalPriceMatched := 0.00
	totalPriceUnknown := 0.00
	totalPrice := 0.00

	for _, t := range transactions {
		if t.GetStatus() == model.Matched {
			nbMatched++
			totalPriceMatched += t.GetAmount()
		} else {
			nbUnknown++
			totalPriceUnknown += t.GetAmount()
		}
		totalPrice += t.GetAmount()
		if internal.MustVerbose {
			fmt.Println(t.ToString())
		}
	}

	summarize := []string{
		fmt.Sprintf("%d Transactions Found.", len(transactions)),
		fmt.Sprintf("%d matched", nbMatched),
		fmt.Sprintf("%d unknown", nbUnknown),
		"",
		fmt.Sprintf("total matched : %.2f $", totalPriceMatched),
		fmt.Sprintf("total unknown : %.2f $", totalPriceUnknown),
		fmt.Sprintf("total : %.2f $", totalPrice),
	}

	longest := slices.MaxFunc(summarize, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})
	border := strings.Repeat("*", len(longest)+5)
	println()
	println(border)
	for _, s := range summarize {
		count := max((len(longest) - len(s)), 0)
		fmt.Println("*", s, strings.Repeat(" ", count), "*")
	}

	println(border)
}

func synopsis() {
	log.Fatalf("Vous devez fournir des arguments de traitement comme le file ou le directory")
}

func logAndFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isBlank(val string) bool {
	return len(val) < 1
}

func filter(transactions []*model.Transaction, f string) []*model.Transaction {
	tmp := make([]*model.Transaction, 0)
	for _, t := range transactions {
		if strings.EqualFold(t.GetVendor(), f) {
			tmp = append(tmp, t)
		}
	}
	return slices.Clip(tmp)
}
