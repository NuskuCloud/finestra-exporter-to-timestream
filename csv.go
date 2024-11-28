package main

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

func parseCSV(csvString string) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(csvString))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func parseCSVRow(records [][]string) {

	if len(records) < 2 {
		fmt.Println("No data found in CSV or missing header row")
		return
	}

	// Get the headers
	headers := records[0]

	// Iterate over the data rows starting from the second row to skip the header
	for _, record := range records[1:] {
		for j := 0; j < len(headers); j += 2 {
			// Determine the headers for the columns to process
			header := headers[j]

			// Check if there's a valid column at current index
			col1 := ""
			if j < len(record) {
				col1 = record[j]
			}

			// Check if there's a valid column at next index
			col2 := ""
			if j+1 < len(record) {
				col2 = record[j+1]
			}

			processColumnPair(header, col1, col2, record)
		}
	}
}

const timeLayout = "2006-01-02 15:04:05Z"

func processColumnPair(header string, value string, dateTime string, record []string) {
	/*
	   header = {string} "Temp_Air_Chamber_Rear Value"
	   value = {string} "4.7800000"
	   dateTime = {string} "2024-11-25 00:01:03Z"
	*/
	if len(dateTime) == 0 {
		return
	}

	measureName := strings.Replace(header, " Value", "", -1)
	timestamp, err := parseTimestamp(dateTime)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return
	}

	err = insertToTimestreamAsync(*TimestreamDatabase, *TimestreamTable, measureName, value, timestamp.Unix())
	if err != nil {
		fmt.Println("Error inserting into Timestream:", err)
		return
	}

	// Placeholder for custom processing logic
	fmt.Printf("Processed: %s, %s, %s\n", timestamp.String(), measureName, value)
}

func parseTimestamp(dateTime string) (time.Time, error) {
	return time.Parse(timeLayout, dateTime)
}

func insertToTimestreamAsync(databaseName string, tableName string, measureName string, measureValue string, timeInSeconds int64) error {
	go func() {
		TimestreamErrCh <- insertIntoTimestream(databaseName, tableName, measureName, measureValue, timeInSeconds)
	}()
	return <-TimestreamErrCh
}
