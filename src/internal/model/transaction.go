package model

import (
	"fmt"
	"strconv"
)

type Transaction struct {
	date          string
	description   string
	amount        float64
	accountNumber string

	Vendor   string
	Category string
	Status   Status
}

type Status string

const (
	Matched Status = "MATCHED"
	Unknown Status = "UNKNOWN"
)

func (t *Transaction) SetDate(v string) {
	t.date = v
}

func (t *Transaction) SetDescription(v string) {
	t.description = v
}

func (t *Transaction) SetAmount(v string) error {
	f, err := strconv.ParseFloat(v, 2)
	if err != nil {
		return err
	}
	t.amount = f
	return nil
}

func (t *Transaction) SetAccountNumber(v string) {
	t.accountNumber = v
}

func (t *Transaction) GetDescription() string {
	return t.description
}

func (t *Transaction) GetAmount() float64 {
	return t.amount
}

func (t *Transaction) ToString() string {
	return fmt.Sprintf("%s \t|\t %50s \t| %8s \t| %5f \t| %10s \t| %15s \t| %6s", t.date, t.description, t.accountNumber, t.amount, t.Vendor, t.Category, t.Status)
}
