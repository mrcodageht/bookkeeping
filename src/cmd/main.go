package main

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strings"

	"mrccommunity.com/bookkeping/src/internal"
	"mrccommunity.com/bookkeping/src/internal/matcher"
	"mrccommunity.com/bookkeping/src/internal/model"
	"mrccommunity.com/bookkeping/src/internal/parser"
)

func main() {

	internal.ReadFlag()
	if len(internal.FileArg) == 0 && len(internal.DirArg) == 0 {
		synopsis()
	}

	if internal.MustVerbose {
		fmt.Printf(">> file : %s\n", internal.FileArg)
		fmt.Printf(">> dir : %s\n", internal.DirArg)
	}
	var transactions []*model.Transaction
	var err error
	var vendors *model.VendorsYAML
	var configs *model.Config

	configs, err = internal.LoadConfig("../bookkeping/configs.yaml")
	logAndFail(err)

	vendors, err = internal.LoadVendor(internal.VendorArg)
	logAndFail(err)

	if internal.MustVerbose {
		fmt.Printf("--- configs:\n%v\n\n", configs)
		fmt.Printf("--- vendors:\n%v\n\n", vendors)
	}

	transactions, err = parser.Parser(internal.FileArg, configs.Output.ErrorFile)
	if err != nil {
		log.Fatal(err)
	}

	matcher.Match(transactions, vendors)
	nbMatched := 0
	nbUnknown := 0
	totalPriceMatched := 0.00
	totalPriceUnknown := 0.00
	totalPrice := 0.00

	for _, t := range transactions {
		if t.Status == model.Matched {
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
