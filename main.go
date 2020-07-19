package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func upload(w http.ResponseWriter, req *http.Request) {
	fmt.Println("here")
	// 2 MB
	req.ParseMultipartForm(2 << 20)

	file, header, err := req.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Pedido incorrecto")
		w.WriteHeader(http.StatusBadRequest)

		logrus.Errorf("Failed to process request: %v\n\n", err)
		return
	}
	defer file.Close()

	fmt.Printf("File name %s\n", header.Filename)

	record, err := extract(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "No se pudo procesar el documento")

		logrus.Errorf("Failed to process file %s: %s\n\n", header.Filename, err)
		return
	}

	recordJson, err := json.Marshal(record)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "No se pudo procesar el documento")

		logrus.Errorf("Failed to marshal record %s: %s\n\n", header.Filename, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=aranduka.json")
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(recordJson))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", upload)

	port := os.Getenv("PORT")
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func extract(file io.Reader) (*record, error) {
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	recordStruct, err := toStruct(records)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert to struct")
	}

	return recordStruct, nil
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
