package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"mrccommunity.com/bookkeping/src/internal"
	"mrccommunity.com/bookkeping/src/internal/model"
)

func Parser(filePath string, outputErr string) ([]*model.Transaction, error) {
	var file *os.File
	var err error
	transactions := make([]*model.Transaction, 0, 5)

	file, err = os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors la lecture du fichier %v", err)
	}
	println("Reading", file.Name(), "...")
	defer file.Close()
	reader := csv.NewReader(file)
	if internal.MustVerbose {
		log.Println("Debut de la lecture du fichier")
	}
	for {
		var record []string
		record, err = reader.Read()

		if err == io.EOF {
			if internal.MustVerbose {
				log.Println("Fin du fichier")
			}
			transactions = slices.Clip(transactions)
			return transactions, nil
		}

		if err != nil {
			return nil, fmt.Errorf("Erreur lors la lecture du fichier %v", err)
		}
		t, err := toTransaction(record)
		if err != nil {
			err = nil
			deadLine(outputErr, record)
			continue
		}
		transactions = append(transactions, t)
	}
}

func toTransaction(record []string) (*model.Transaction, error) {
	transaction := model.Transaction{}
	var err error
	for index, field := range record {

		switch index {
		case 0, 3, 7, 5:
			continue
		case 1:
			transaction.SetAccountNumber(field)
		case 2:
			transaction.SetDate(field)
		case 4:
			transaction.SetDescription(field)
		case 6:
			err = transaction.SetAmount(field)
		}

	}
	return &transaction, err
}

func deadLine(output string, record []string) error {
	file, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if _, ok := err.(*os.PathError); ok {
		file, err = os.Create(output)
	}
	if err != nil {
		return fmt.Errorf("Erreur lors la lecture du fichier %v", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.Write(record)
	if err != nil {
		fmt.Printf("Erreur lors de l'ecriture du fichier %v", err)
	}
	writer.Flush()
	if err = writer.Error(); err != nil {
		fmt.Println(err)
	}
	return nil
}
