package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func main() {
	flFile := flag.String("file", "", "the file to convert")
	flag.Parse()

	recordFile, err := os.Open(*flFile)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}

	reader := csv.NewReader(recordFile)
	records, _ := reader.ReadAll()
	recordStruct, err := toStruct(records)
	if err != nil {
		log.Fatalf("failed to convert to struct: %v", err)
	}

	recordJson, err := json.Marshal(recordStruct)
	if err != nil {
		log.Fatalf("failed to marshal json: %v", err)
	}

	file, err := os.Create("out.json")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	file.WriteString(string(recordJson))
	file.Close()
}

func toStruct(records [][]string) (*record, error) {

	var expenses []*expense
	for i := 1; i < len(records); i++ {
		record := records[i]

		date, err := time.Parse("2006-01-02", record[2])
		if err != nil {
			return nil, errors.Wrap(err, "failed to get date")
		}
		period := strconv.Itoa(date.Year())

		amount, err := strconv.ParseInt(record[9], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "invalid amount")
		}

		expenses = append(expenses, &expense{
			Period:            period,
			RUC:               record[0],
			DocumentType:      documentTypeID(record[1]),
			DocumentTypeText:  record[1],
			Date:              record[2],
			TimbradoNumber:    record[3],
			TimbradoDocument:  record[4],
			TimbradoCondition: record[5],
			EntityIDType:      record[6],
			EntityID:          record[7],
			Entity:            record[8],
			Amount:            amount,
			ExpenseType:       expenseType(record[10]),
			ExpenseSubtype:    expenseSubtype(record[11]),
		})
	}

	return &record{Expenses: expenses, Incomes: []*income{}}, nil
}
